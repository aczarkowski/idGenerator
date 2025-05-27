package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"uidGenerator/timeprovider/epoch"
	"uidGenerator/timeprovider/julian"

	"github.com/labstack/echo/v4"
)

func TestMain_Integration(t *testing.T) {
	// Test that main components work together
	// This is a basic integration test without starting the full server
	
	// Test epoch time provider integration
	epochProvider := epoch.New(1420070400000)
	if epochProvider == nil {
		t.Error("Failed to create epoch provider")
	}
	
	// Test julian time provider integration
	julianProvider := julian.New(2000100000)
	if julianProvider == nil {
		t.Error("Failed to create julian provider")
	}
	
	// Test that both providers return reasonable timestamps
	epochTime := epochProvider.GetTimeStamp()
	julianTime := julianProvider.GetTimeStamp()
	
	if epochTime <= 0 {
		t.Errorf("Epoch provider returned invalid timestamp: %d", epochTime)
	}
	
	if julianTime <= 0 {
		t.Errorf("Julian provider returned invalid timestamp: %d", julianTime)
	}
}

func TestFullStack_EpochProvider(t *testing.T) {
	// Create Echo instance
	e := echo.New()
	
	// Setup middleware
	provider := epoch.New(1420070400000)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Simplified version of the middleware for testing
			worker := &struct {
				WorkerID     int64
				ThreadId     int64
				TimeProvider interface{ GetTimeStamp() int64 }
			}{
				WorkerID:     1,
				ThreadId:     1,
				TimeProvider: provider,
			}
			c.Set("worker", worker)
			return next(c)
		}
	})
	
	// Add a simple route for testing
	e.GET("/test", func(c echo.Context) error {
		worker := c.Get("worker")
		if worker == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "no worker"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	
	// Test the route
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", response["status"])
	}
}

func TestFullStack_JulianProvider(t *testing.T) {
	// Create Echo instance
	e := echo.New()
	
	// Setup middleware with Julian provider
	provider := julian.New(2000100000)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Simplified version of the middleware for testing
			worker := &struct {
				WorkerID     int64
				ThreadId     int64
				TimeProvider interface{ GetTimeStamp() int64 }
			}{
				WorkerID:     1,
				ThreadId:     1,
				TimeProvider: provider,
			}
			c.Set("worker", worker)
			return next(c)
		}
	})
	
	// Add a simple route for testing
	e.GET("/test", func(c echo.Context) error {
		worker := c.Get("worker")
		if worker == nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "no worker"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})
	
	// Test the route
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	var response map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", response["status"])
	}
}

func TestTimeProviders_Comparison(t *testing.T) {
	// Test that different time providers give different results
	epochProvider := epoch.New(1420070400000)
	julianProvider := julian.New(2000100000)
	
	epochTime1 := epochProvider.GetTimeStamp()
	julianTime1 := julianProvider.GetTimeStamp()
	
	// They should be different (different algorithms)
	if epochTime1 == julianTime1 {
		t.Log("Epoch and Julian timestamps are the same, which is unlikely but possible")
	}
	
	// Both should be positive
	if epochTime1 <= 0 || julianTime1 <= 0 {
		t.Errorf("Expected positive timestamps, got epoch: %d, julian: %d", epochTime1, julianTime1)
	}
	
	// Test temporal ordering (later calls should give >= timestamps)
	time.Sleep(time.Millisecond * 2)
	
	epochTime2 := epochProvider.GetTimeStamp()
	julianTime2 := julianProvider.GetTimeStamp()
	
	if epochTime2 < epochTime1 {
		t.Errorf("Expected epoch time to increase, got %d then %d", epochTime1, epochTime2)
	}
	
	if julianTime2 < julianTime1 {
		t.Errorf("Expected julian time to increase, got %d then %d", julianTime1, julianTime2)
	}
}

func TestConfiguration_DefaultValues(t *testing.T) {
	// Test that the default configuration values are reasonable
	// This tests the constants defined in main.go conceptually
	
	// Test default offset for epoch (Jan 1, 2015)
	defaultOffset := int64(1420070400000)
	provider := epoch.New(defaultOffset)
	timestamp := provider.GetTimeStamp()
	
	// Should be positive after subtracting the offset
	if timestamp <= 0 {
		t.Errorf("Expected positive timestamp with default offset, got %d", timestamp)
	}
	
	// Test default offset for julian
	julianOffset := int64(2000100000)
	julianProvider := julian.New(julianOffset)
	julianTimestamp := julianProvider.GetTimeStamp()
	
	// Should be positive after subtracting the offset
	if julianTimestamp <= 0 {
		t.Errorf("Expected positive julian timestamp with default offset, got %d", julianTimestamp)
	}
}
