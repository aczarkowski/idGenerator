package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"uidGenerator/generator"
	"uidGenerator/timeprovider/epoch"

	"github.com/labstack/echo/v4"
)

func TestGenerator_SingleID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup worker in context
	provider := epoch.New(1420070400000)
	worker := &generator.WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	c.Set("worker", worker)

	// Call handler
	if err := Generator(c); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check response
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	ids, ok := response["ids"].([]interface{})
	if !ok {
		t.Error("Expected 'ids' field in response")
	}

	if len(ids) != 1 {
		t.Errorf("Expected 1 ID, got %d", len(ids))
	}
}

func TestGenerator_MultipleIDs(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?numberOfIds=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup worker in context
	provider := epoch.New(1420070400000)
	worker := &generator.WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	c.Set("worker", worker)

	// Call handler
	if err := Generator(c); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check response
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	ids, ok := response["ids"].([]interface{})
	if !ok {
		t.Error("Expected 'ids' field in response")
	}

	if len(ids) != 5 {
		t.Errorf("Expected 5 IDs, got %d", len(ids))
	}
}

func TestGenerator_InvalidNumberOfIds(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?numberOfIds=invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup worker in context
	provider := epoch.New(1420070400000)
	worker := &generator.WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	c.Set("worker", worker)

	// Call handler
	if err := Generator(c); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Should default to 1 ID when parsing fails
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	ids, ok := response["ids"].([]interface{})
	if !ok {
		t.Error("Expected 'ids' field in response")
	}

	if len(ids) != 1 {
		t.Errorf("Expected 1 ID (default), got %d", len(ids))
	}
}

func TestGenerator_ResponseFormat(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Setup worker in context
	provider := epoch.New(1420070400000)
	worker := &generator.WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	c.Set("worker", worker)

	// Call handler
	if err := Generator(c); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Check content type
	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got %s", contentType)
	}

	// Check that response is valid JSON
	var response map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Check response structure
	if _, exists := response["ids"]; !exists {
		t.Error("Response should contain 'ids' field")
	}
}
