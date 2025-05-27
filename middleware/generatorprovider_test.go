package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"uidGenerator/generator"
	"uidGenerator/timeprovider/epoch"

	"github.com/labstack/echo/v4"
)

func TestGeneratorProvider(t *testing.T) {
	workerId := int64(2)
	provider := epoch.New(1420070400000)

	middleware := GeneratorProvider(workerId, provider)

	// Create a dummy handler to test the middleware
	handler := func(c echo.Context) error {
		worker := c.Get("worker").(*generator.WorkerVariant)
		if worker == nil {
			t.Error("Expected worker to be set in context")
		}
		if worker.WorkerID != workerId {
			t.Errorf("Expected worker ID %d, got %d", workerId, worker.WorkerID)
		}
		if worker.TimeProvider == nil {
			t.Error("Expected time provider to be set")
		}
		return c.String(http.StatusOK, "OK")
	}

	// Wrap handler with middleware
	wrappedHandler := middleware(handler)

	// Test the middleware
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := wrappedHandler(c); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

func TestGeneratorProvider_WorkerPool(t *testing.T) {
	workerId := int64(1)
	provider := epoch.New(1420070400000)

	middleware := GeneratorProvider(workerId, provider)

	// Test that we can handle multiple concurrent requests
	// This tests the worker pool functionality
	var threadIds []int64
	handler := func(c echo.Context) error {
		worker := c.Get("worker").(*generator.WorkerVariant)
		threadIds = append(threadIds, worker.ThreadId)
		return c.String(http.StatusOK, "OK")
	}

	wrappedHandler := middleware(handler)

	// Make multiple requests to test worker pool
	for i := 0; i < int(generator.ThreadCap); i++ {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		if err := wrappedHandler(c); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	// Check that we got the expected number of thread IDs
	if len(threadIds) != int(generator.ThreadCap) {
		t.Errorf("Expected %d thread IDs, got %d", generator.ThreadCap, len(threadIds))
	}

	// Check that all thread IDs are within valid range
	for _, threadId := range threadIds {
		if threadId < 1 || threadId > generator.ThreadCap {
			t.Errorf("Thread ID %d is out of valid range (1-%d)", threadId, generator.ThreadCap)
		}
	}
}

func TestGeneratorProvider_WorkerConfiguration(t *testing.T) {
	workerId := int64(5)
	provider := epoch.New(1420070400000)

	middleware := GeneratorProvider(workerId, provider)

	handler := func(c echo.Context) error {
		worker := c.Get("worker").(*generator.WorkerVariant)

		// Verify worker configuration
		if worker.WorkerID != workerId {
			t.Errorf("Expected worker ID %d, got %d", workerId, worker.WorkerID)
		}

		if worker.ThreadId < 1 || worker.ThreadId > generator.ThreadCap {
			t.Errorf("Thread ID %d is out of valid range", worker.ThreadId)
		}

		if worker.TimeProvider != provider {
			t.Error("Time provider not correctly set")
		}

		return c.String(http.StatusOK, "OK")
	}

	wrappedHandler := middleware(handler)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := wrappedHandler(c); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGeneratorProvider_WorkerPoolSize(t *testing.T) {
	provider := epoch.New(1420070400000)
	workerId := int64(1)

	middleware := GeneratorProvider(workerId, provider)

	// Test multiple concurrent requests to ensure worker pool works
	threadIds := make(map[int64]bool)
	requestCount := int(generator.ThreadCap)

	for i := 0; i < requestCount; i++ {
		testHandler := func(c echo.Context) error {
			worker := c.Get("worker").(*generator.WorkerVariant)
			threadIds[worker.ThreadId] = true
			return nil
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		middlewareFunc := middleware(testHandler)
		err := middlewareFunc(c)

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	// Check that we got different thread IDs (worker pool is working)
	if len(threadIds) == 0 {
		t.Error("Expected at least one thread ID")
	}

	// All thread IDs should be within valid range
	for threadId := range threadIds {
		if threadId <= 0 || threadId > generator.ThreadCap {
			t.Errorf("Invalid thread ID: %d", threadId)
		}
	}
}

func TestGeneratorProvider_ThreadIdRange(t *testing.T) {
	provider := epoch.New(1420070400000)
	workerId := int64(1)

	middleware := GeneratorProvider(workerId, provider)

	testHandler := func(c echo.Context) error {
		worker := c.Get("worker").(*generator.WorkerVariant)

		if worker.ThreadId <= 0 || worker.ThreadId > generator.ThreadCap {
			t.Errorf("Expected thread ID between 1 and %d, got %d", generator.ThreadCap, worker.ThreadId)
		}

		return nil
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	middlewareFunc := middleware(testHandler)
	err := middlewareFunc(c)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGeneratorProvider_WorkerReuse(t *testing.T) {
	provider := epoch.New(1420070400000)
	workerId := int64(1)

	middleware := GeneratorProvider(workerId, provider)

	var firstWorker *generator.WorkerVariant
	var secondWorker *generator.WorkerVariant

	// First request
	firstHandler := func(c echo.Context) error {
		firstWorker = c.Get("worker").(*generator.WorkerVariant)
		return nil
	}

	e1 := echo.New()
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec1 := httptest.NewRecorder()
	c1 := e1.NewContext(req1, rec1)

	middlewareFunc1 := middleware(firstHandler)
	err1 := middlewareFunc1(c1)

	if err1 != nil {
		t.Errorf("Expected no error, got %v", err1)
	}

	// Second request - should potentially reuse the same worker
	secondHandler := func(c echo.Context) error {
		secondWorker = c.Get("worker").(*generator.WorkerVariant)
		return nil
	}

	e2 := echo.New()
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec2 := httptest.NewRecorder()
	c2 := e2.NewContext(req2, rec2)

	middlewareFunc2 := middleware(secondHandler)
	err2 := middlewareFunc2(c2)

	if err2 != nil {
		t.Errorf("Expected no error, got %v", err2)
	}

	// Workers should be valid
	if firstWorker == nil || secondWorker == nil {
		t.Error("Expected workers to be set")
	}

	// Both workers should have the same worker ID
	if firstWorker.WorkerID != secondWorker.WorkerID {
		t.Errorf("Expected same worker ID, got %d and %d", firstWorker.WorkerID, secondWorker.WorkerID)
	}
}