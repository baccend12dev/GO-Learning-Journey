package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend/config"
	"backend/models"
	"backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupFeatureRequestTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate related schemas
	err = db.AutoMigrate(&models.Server{}, &models.System{}, &models.FeatureRequest{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	config.DB = db
}

func TestGetFeatureRequestsBySystemID(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system
	system := models.System{Name: "Core Portal", Type: "Web", Links: "http://core", Status: "Active"}
	config.DB.Create(&system)

	// Seed feature requests
	requests := []models.FeatureRequest{
		{SystemId: system.ID, Title: "Dark Mode", Description: "Implement system wide dark mode", Status: "Pending"},
		{SystemId: system.ID, Title: "Export PDF", Description: "Export page data to PDF", Status: "In Progress"},
	}
	for _, req := range requests {
		config.DB.Create(&req)
	}

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	req, _ := http.NewRequest(http.MethodGet, "/api/systems/"+strconv.Itoa(int(system.ID))+"/feature-requests", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var returnedRequests []models.FeatureRequest
	if err := json.Unmarshal(w.Body.Bytes(), &returnedRequests); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if len(returnedRequests) != 2 {
		t.Errorf("Expected 2 feature requests, but got %d", len(returnedRequests))
	}

	if returnedRequests[0].Title != "Dark Mode" || returnedRequests[1].Title != "Export PDF" {
		t.Errorf("Returned feature request data does not match seeded records")
	}
}

func TestCreateFeatureRequest(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system
	system := models.System{Name: "Core Portal", Type: "Web", Links: "http://core", Status: "Active"}
	config.DB.Create(&system)

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	// Case 1: Success Creation
	reqBody := models.CreateFeatureRequest{
		Title:       "User Management",
		Description: "Add complete role based access control",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/feature-requests", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	var responseReq models.FeatureRequest
	if err := json.Unmarshal(w.Body.Bytes(), &responseReq); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseReq.Title != reqBody.Title || responseReq.Description != reqBody.Description {
		t.Errorf("Response fields mismatch with request body")
	}

	if responseReq.Status != "Pending" {
		t.Errorf("Expected default status to be 'Pending', but got %q", responseReq.Status)
	}

	// Verify in database
	var dbReq models.FeatureRequest
	if err := config.DB.First(&dbReq, responseReq.ID).Error; err != nil {
		t.Errorf("Feature Request was not found in the database: %v", err)
	}

	// Case 2: System not found
	wNotFound := httptest.NewRecorder()
	reqNotFound, _ := http.NewRequest(http.MethodPost, "/api/systems/999/feature-requests", bytes.NewBuffer(bodyBytes))
	reqNotFound.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d when system is missing, but got %d", http.StatusNotFound, wNotFound.Code)
	}

	// Case 3: Invalid request body (missing description)
	invalidReq := map[string]interface{}{
		"title": "Invalid Request",
	}
	invalidBytes, _ := json.Marshal(invalidReq)
	wInvalid := httptest.NewRecorder()
	reqInvalid, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/feature-requests", bytes.NewBuffer(invalidBytes))
	reqInvalid.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wInvalid, reqInvalid)

	if wInvalid.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d on validation failure, but got %d", http.StatusBadRequest, wInvalid.Code)
	}
}

func TestUpdateFeatureRequest(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system & feature request
	system := models.System{Name: "Core Portal"}
	config.DB.Create(&system)

	feature := models.FeatureRequest{SystemId: system.ID, Title: "Old Title", Description: "Old Desc", Status: "Pending"}
	config.DB.Create(&feature)

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	// Case 1: Success update status and details
	updatePayload := map[string]interface{}{
		"title":       "New Title",
		"description": "New Description",
		"status":      "Completed",
	}
	bodyBytes, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, "/api/feature-requests/"+strconv.Itoa(int(feature.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var returnedFeature models.FeatureRequest
	if err := json.Unmarshal(w.Body.Bytes(), &returnedFeature); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if returnedFeature.Title != "New Title" || returnedFeature.Status != "Completed" {
		t.Errorf("Returned object does not match updated values")
	}

	// Verify in database
	var dbFeature models.FeatureRequest
	config.DB.First(&dbFeature, feature.ID)
	if dbFeature.Title != "New Title" || dbFeature.Status != "Completed" {
		t.Errorf("Database record was not updated correctly")
	}

	// Case 2: Not Found Update
	wNotFound := httptest.NewRecorder()
	reqNotFound, _ := http.NewRequest(http.MethodPut, "/api/feature-requests/999", bytes.NewBuffer(bodyBytes))
	reqNotFound.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d on updating non-existent request, but got %d", http.StatusNotFound, wNotFound.Code)
	}
}

func TestDeleteFeatureRequest(t *testing.T) {
	setupFeatureRequestTestDB(t)

	// Seed system & feature request
	system := models.System{Name: "Core Portal"}
	config.DB.Create(&system)

	feature := models.FeatureRequest{SystemId: system.ID, Title: "To Delete", Description: "Desc", Status: "Pending"}
	config.DB.Create(&feature)

	r := gin.Default()
	routes.SetupFeatureRequestRoutes(r)

	// Case 1: Success deletion
	req, _ := http.NewRequest(http.MethodDelete, "/api/feature-requests/"+strconv.Itoa(int(feature.ID)), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	// Verify database record is deleted
	var dbFeature models.FeatureRequest
	err := config.DB.First(&dbFeature, feature.ID).Error
	if err == nil {
		t.Errorf("Feature request record still exists in the database")
	}

	// Case 2: Not Found Deletion
	wNotFound := httptest.NewRecorder()
	reqNotFound, _ := http.NewRequest(http.MethodDelete, "/api/feature-requests/999", nil)
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d on deleting non-existent request, but got %d", http.StatusNotFound, wNotFound.Code)
	}
}
