// +build e2e

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"module_6/cmd/server/app"
	"module_6/internal/models"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestApp(t *testing.T) *fiber.App {
	// Set test environment variables
	t.Setenv("ENV", "development")
	t.Setenv("DB_NAME", "notes_db_test")
	t.Setenv("JWT_SECRET", "test-secret-key")
	t.Setenv("LOG_LEVEL", "error") // Reduce noise in tests
	t.Setenv("LOG_OUTPUT", "stdout")

	application, err := app.New()
	require.NoError(t, err)
	
	return application.GetFiberApp() // You'll need to add this method to app
}

func TestE2E_CompleteUserFlow(t *testing.T) {
	app := setupTestApp(t)

	var token string
	var categoryID int
	var noteID int

	t.Run("1. User Registration", func(t *testing.T) {
		payload := map[string]string{
			"email":    "e2e@example.com",
			"name":     "e2euser",
			"password": "password123",
		}
		
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
		
		var result map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		
		token = result["token"].(string)
		assert.NotEmpty(t, token)
	})

	t.Run("2. Login with same credentials", func(t *testing.T) {
		payload := map[string]string{
			"email":    "e2e@example.com",
			"password": "password123",
		}
		
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/signin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("3. Create Category", func(t *testing.T) {
		payload := map[string]string{
			"name": "E2E Test Category",
		}
		
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
		
		var result map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		
		categoryID = int(result["id"].(float64))
		assert.NotZero(t, categoryID)
	})

	t.Run("4. Get Categories", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/categories", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		
		var result map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		
		categories := result["categories"].([]interface{})
		assert.Len(t, categories, 1)
	})

	t.Run("5. Create Note", func(t *testing.T) {
		payload := map[string]interface{}{
			"title":       "E2E Test Note",
			"content":     "This is a test note content",
			"category_id": categoryID,
		}
		
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/notes", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
		
		var result map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		
		noteID = int(result["id"].(float64))
		assert.NotZero(t, noteID)
	})

	t.Run("6. Get Notes", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/notes", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
		
		var result map[string]interface{}
		respBody, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBody, &result)
		
		notes := result["notes"].([]interface{})
		assert.Len(t, notes, 1)
	})

	t.Run("7. Update Note", func(t *testing.T) {
		payload := map[string]string{
			"title":   "Updated E2E Test Note",
			"content": "Updated content",
		}
		
		body, _ := json.Marshal(payload)
		url := fmt.Sprintf("/api/notes/%d", noteID)
		req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("8. Delete Note", func(t *testing.T) {
		url := fmt.Sprintf("/api/notes/%d", noteID)
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 204, resp.StatusCode)
	})

	t.Run("9. Delete Category", func(t *testing.T) {
		url := fmt.Sprintf("/api/categories/%d", categoryID)
		req := httptest.NewRequest("DELETE", url, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestE2E_ErrorCases(t *testing.T) {
	app := setupTestApp(t)

	t.Run("Invalid credentials", func(t *testing.T) {
		payload := map[string]string{
			"email":    "wrong@example.com",
			"password": "wrongpassword",
		}
		
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", "/api/auth/signin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("Access without token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/categories", nil)
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 401, resp.StatusCode)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/auth/signup", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
}