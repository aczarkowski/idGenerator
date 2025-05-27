package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"uidGenerator/timeprovider/epoch"
	"uidGenerator/timeprovider/julian"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"uidGenerator/handler"
	generatorMiddleware "uidGenerator/middleware"
)

func TestIntegration_EpochProvider(t *testing.T) {
	// Setup server with epoch provider
	e := echo.New()
	provider := epoch.New(1420070400000)
	workerId := int64(1)
	
	e.Use(generatorMiddleware.GeneratorProvider(workerId, provider))
	e.Use(middleware.Logger())
	e.GET("/", handler.Generator)
	
	// Test single ID generation
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	ids, ok := response["ids"].([]interface{})
	if !ok || len(ids) != 1 {
		t.Error("Expected single ID in response")
	}
}

func TestIntegration_JulianProvider(t *testing.T) {
	// Setup server with julian provider
	e := echo.New()
	provider := julian.New(2000100000)
	workerId := int64(2)
	
	e.Use(generatorMiddleware.GeneratorProvider(workerId, provider))
	e.Use(middleware.Logger())
	e.GET("/", handler.Generator)
	
	// Test multiple ID generation
	req := httptest.NewRequest(http.MethodGet, "/?numberOfIds=3", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	ids, ok := response["ids"].([]interface{})
	if !ok || len(ids) != 3 {
		t.Errorf("Expected 3 IDs in response, got %d", len(ids))
	}
}

func TestIntegration_MultipleRequests(t *testing.T) {
	// Setup server
	e := echo.New()
	provider := epoch.New(1420070400000)
	workerId := int64(3)
	
	e.Use(generatorMiddleware.GeneratorProvider(workerId, provider))
	e.GET("/", handler.Generator)
	
	allIds := make(map[int64]bool)
	
	// Make multiple requests and ensure all IDs are unique
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest(http.MethodGet, "/?numberOfIds=5", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Errorf("Request %d: Expected status 200, got %d", i, rec.Code)
			continue
		}
		
		var response map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Errorf("Request %d: Failed to parse response: %v", i, err)
			continue
		}
		
		ids, ok := response["ids"].([]interface{})
		if !ok {
			t.Errorf("Request %d: No ids in response", i)
			continue
		}
		
		for _, idInterface := range ids {
			// Convert to int64 (JSON numbers come as float64)
			idFloat, ok := idInterface.(float64)
			if !ok {
				t.Errorf("Request %d: ID is not a number: %v", i, idInterface)
				continue
			}
			id := int64(idFloat)
			
			if allIds[id] {
				t.Errorf("Request %d: Duplicate ID found: %d", i, id)
			}
			allIds[id] = true
		}
	}
	
	// Should have 50 unique IDs (10 requests * 5 IDs each)
	if len(allIds) != 50 {
		t.Errorf("Expected 50 unique IDs, got %d", len(allIds))
	}
}

func TestIntegration_LargeNumberOfIds(t *testing.T) {
	// Setup server
	e := echo.New()
	provider := epoch.New(1420070400000)
	workerId := int64(4)
	
	e.Use(generatorMiddleware.GeneratorProvider(workerId, provider))
	e.GET("/", handler.Generator)
	
	// Request 100 IDs
	req := httptest.NewRequest(http.MethodGet, "/?numberOfIds=100", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}
	
	ids, ok := response["ids"].([]interface{})
	if !ok || len(ids) != 100 {
		t.Errorf("Expected 100 IDs in response, got %d", len(ids))
	}
	
	// Verify all IDs are unique
	idMap := make(map[int64]bool)
	for _, idInterface := range ids {
		idFloat := idInterface.(float64)
		id := int64(idFloat)
		if idMap[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		idMap[id] = true
	}
}

func TestIntegration_DifferentWorkerIds(t *testing.T) {
	// Test that different worker IDs produce different ID patterns
	provider := epoch.New(1420070400000)
	
	// Worker 1
	e1 := echo.New()
	e1.Use(generatorMiddleware.GeneratorProvider(1, provider))
	e1.GET("/", handler.Generator)
	
	// Worker 2
	e2 := echo.New()
	e2.Use(generatorMiddleware.GeneratorProvider(2, provider))
	e2.GET("/", handler.Generator)
	
	// Get IDs from both workers
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec1 := httptest.NewRecorder()
	e1.ServeHTTP(rec1, req1)
	
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	rec2 := httptest.NewRecorder()
	e2.ServeHTTP(rec2, req2)
	
	// Parse responses
	var response1, response2 map[string]interface{}
	json.Unmarshal(rec1.Body.Bytes(), &response1)
	json.Unmarshal(rec2.Body.Bytes(), &response2)
	
	ids1 := response1["ids"].([]interface{})
	ids2 := response2["ids"].([]interface{})
	
	id1 := int64(ids1[0].(float64))
	id2 := int64(ids2[0].(float64))
	
	// IDs should be different (very high probability with different worker IDs)
	if id1 == id2 {
		t.Error("Expected different IDs from different workers, but got the same")
	}
}

func TestIntegration_InvalidParameters(t *testing.T) {
	// Setup server
	e := echo.New()
	provider := epoch.New(1420070400000)
	workerId := int64(5)
	
	e.Use(generatorMiddleware.GeneratorProvider(workerId, provider))
	e.GET("/", handler.Generator)
	
	testCases := []struct {
		query          string
		expectedIds    int
		description    string
	}{
		{"?numberOfIds=abc", 1, "non-numeric parameter"},
		{"?numberOfIds=-5", 1, "negative parameter"},
		{"?numberOfIds=0", 1, "zero parameter"},
		{"?numberOfIds=", 1, "empty parameter"},
		{"?invalidParam=5", 1, "invalid parameter name"},
	}
	
	for _, tc := range testCases {
		req := httptest.NewRequest(http.MethodGet, "/"+tc.query, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		
		if rec.Code != http.StatusOK {
			t.Errorf("%s: Expected status 200, got %d", tc.description, rec.Code)
			continue
		}
		
		var response map[string]interface{}
		if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
			t.Errorf("%s: Failed to parse response: %v", tc.description, err)
			continue
		}
		
		ids, ok := response["ids"].([]interface{})
		if !ok || len(ids) != tc.expectedIds {
			t.Errorf("%s: Expected %d IDs, got %d", tc.description, tc.expectedIds, len(ids))
		}
	}
}
