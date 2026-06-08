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

func setupDocumentationTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Server{}, &models.System{}, &models.Documentation{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	config.DB = db
}

func TestGetDocumentationsBySystemID(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	docs := []models.Documentation{
		{SystemId: system.ID, Title: "Flow Bisnis Cuti", Category: "Business Flow", Content: "Diagram flow cuti karyawan"},
		{SystemId: system.ID, Title: "Panduan Deploy", Category: "Deployment Guide", Content: "Langkah build & deploy"},
	}
	for _, doc := range docs {
		config.DB.Create(&doc)
	}

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	// Test case 1: List all for system
	req, _ := http.NewRequest(http.MethodGet, "/api/systems/"+strconv.Itoa(int(system.ID))+"/documentations", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var returnedDocs []models.Documentation
	json.Unmarshal(w.Body.Bytes(), &returnedDocs)
	if len(returnedDocs) != 2 {
		t.Errorf("Expected 2 documentations, got %d", len(returnedDocs))
	}

	// Test case 2: List filtered by category
	reqFilter, _ := http.NewRequest(http.MethodGet, "/api/systems/"+strconv.Itoa(int(system.ID))+"/documentations?category=Deployment Guide", nil)
	wFilter := httptest.NewRecorder()
	r.ServeHTTP(wFilter, reqFilter)

	if wFilter.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", wFilter.Code)
	}

	var filteredDocs []models.Documentation
	json.Unmarshal(wFilter.Body.Bytes(), &filteredDocs)
	if len(filteredDocs) != 1 || filteredDocs[0].Title != "Panduan Deploy" {
		t.Errorf("Expected 1 filtered documentation (Panduan Deploy), got: %v", filteredDocs)
	}
}

func TestGetDocumentationByID(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	doc := models.Documentation{SystemId: system.ID, Title: "Database Schema", Category: "Database Documentation", Content: "Table list"}
	config.DB.Create(&doc)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	// Case 1: Success fetch
	req, _ := http.NewRequest(http.MethodGet, "/api/documentations/"+strconv.Itoa(int(doc.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var returnedDoc models.Documentation
	json.Unmarshal(w.Body.Bytes(), &returnedDoc)
	if returnedDoc.Title != doc.Title {
		t.Errorf("Expected title %q, got %q", doc.Title, returnedDoc.Title)
	}

	// Case 2: Not found
	reqNotFound, _ := http.NewRequest(http.MethodGet, "/api/documentations/999", nil)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", wNotFound.Code)
	}
}

func TestCreateDocumentation(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	// Case 1: Success creation
	payload := models.CreateDocumentationRequest{
		Title:    "API List",
		Category: "API Documentation",
		Content:  "List of REST APIs",
	}
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/documentations", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	// Case 2: Invalid category validation
	invalidPayload := models.CreateDocumentationRequest{
		Title:    "API List",
		Category: "Random Category Name",
		Content:  "Content",
	}
	invalidBytes, _ := json.Marshal(invalidPayload)
	reqInvalid, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/documentations", bytes.NewBuffer(invalidBytes))
	reqInvalid.Header.Set("Content-Type", "application/json")
	wInvalid := httptest.NewRecorder()
	r.ServeHTTP(wInvalid, reqInvalid)

	if wInvalid.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid category validation, got %d", wInvalid.Code)
	}

	// Case 3: System not found
	reqSystemNotFound, _ := http.NewRequest(http.MethodPost, "/api/systems/999/documentations", bytes.NewBuffer(bodyBytes))
	reqSystemNotFound.Header.Set("Content-Type", "application/json")
	wSystemNotFound := httptest.NewRecorder()
	r.ServeHTTP(wSystemNotFound, reqSystemNotFound)

	if wSystemNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 when system is missing, got %d", wSystemNotFound.Code)
	}
}

func TestUpdateDocumentation(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	doc := models.Documentation{SystemId: system.ID, Title: "Draft Panduan", Category: "User Manual", Content: "Draft content"}
	config.DB.Create(&doc)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	// Case 1: Success update
	updatePayload := map[string]interface{}{
		"title":    "Panduan Penggunaan Akhir",
		"category": "User Manual",
		"content":  "Final content",
	}
	bodyBytes, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, "/api/documentations/"+strconv.Itoa(int(doc.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var updatedDoc models.Documentation
	json.Unmarshal(w.Body.Bytes(), &updatedDoc)
	if updatedDoc.Title != "Panduan Penggunaan Akhir" || updatedDoc.Content != "Final content" {
		t.Errorf("Expected updated fields, got: %v", updatedDoc)
	}

	// Case 2: Update with invalid category
	invalidPayload := map[string]interface{}{
		"category": "Invalid Cat",
	}
	invalidBytes, _ := json.Marshal(invalidPayload)
	reqInvalid, _ := http.NewRequest(http.MethodPut, "/api/documentations/"+strconv.Itoa(int(doc.ID)), bytes.NewBuffer(invalidBytes))
	reqInvalid.Header.Set("Content-Type", "application/json")
	wInvalid := httptest.NewRecorder()
	r.ServeHTTP(wInvalid, reqInvalid)

	if wInvalid.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for updating with invalid category, got %d", wInvalid.Code)
	}
}

func TestDeleteDocumentation(t *testing.T) {
	setupDocumentationTestDB(t)

	system := models.System{Name: "Portal HRD"}
	config.DB.Create(&system)

	doc := models.Documentation{SystemId: system.ID, Title: "To Delete", Category: "Technical Flow", Content: "Temp documentation"}
	config.DB.Create(&doc)

	r := gin.Default()
	routes.SetupDocumentationRoutes(r)

	// Case 1: Success deletion
	req, _ := http.NewRequest(http.MethodDelete, "/api/documentations/"+strconv.Itoa(int(doc.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var dbDoc models.Documentation
	err := config.DB.First(&dbDoc, doc.ID).Error
	if err == nil {
		t.Errorf("Documentation was not deleted from database")
	}

	// Case 2: Not found deletion
	reqNotFound, _ := http.NewRequest(http.MethodDelete, "/api/documentations/999", nil)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", wNotFound.Code)
	}
}
