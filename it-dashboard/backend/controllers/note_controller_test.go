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

func setupNoteTestDB(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	err = db.AutoMigrate(&models.Server{}, &models.System{}, &models.Note{})
	if err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	config.DB = db
}

func TestGetNotesBySystemID(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	notes := []models.Note{
		{SystemId: system.ID, Title: "Note 1", Content: "Content 1"},
		{SystemId: system.ID, Title: "Note 2", Content: "Content 2"},
	}
	for _, note := range notes {
		config.DB.Create(&note)
	}

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	req, _ := http.NewRequest(http.MethodGet, "/api/systems/"+strconv.Itoa(int(system.ID))+"/notes", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	var returnedNotes []models.Note
	if err := json.Unmarshal(w.Body.Bytes(), &returnedNotes); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	if len(returnedNotes) != 2 {
		t.Errorf("Expected 2 notes, got %d", len(returnedNotes))
	}
}

func TestCreateNote(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	// Case 1: Success creation
	payload := models.CreateNoteRequest{
		Title:   "New Note",
		Content: "New Content",
	}
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/api/systems/"+strconv.Itoa(int(system.ID))+"/notes", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	// Case 2: System not found
	reqNotFound, _ := http.NewRequest(http.MethodPost, "/api/systems/999/notes", bytes.NewBuffer(bodyBytes))
	reqNotFound.Header.Set("Content-Type", "application/json")
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", wNotFound.Code)
	}
}

func TestUpdateNote(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	note := models.Note{SystemId: system.ID, Title: "Old Title", Content: "Old Content"}
	config.DB.Create(&note)

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	// Case 1: Success update
	payload := models.CreateNoteRequest{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	bodyBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPut, "/api/notes/"+strconv.Itoa(int(note.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var updatedNote models.Note
	if err := json.Unmarshal(w.Body.Bytes(), &updatedNote); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedNote.Title != payload.Title {
		t.Errorf("Expected updated title %q, got %q", payload.Title, updatedNote.Title)
	}

	// Case 2: Not found
	reqNotFound, _ := http.NewRequest(http.MethodPut, "/api/notes/999", bytes.NewBuffer(bodyBytes))
	reqNotFound.Header.Set("Content-Type", "application/json")
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", wNotFound.Code)
	}
}

func TestDeleteNote(t *testing.T) {
	setupNoteTestDB(t)

	system := models.System{Name: "System A"}
	config.DB.Create(&system)

	note := models.Note{SystemId: system.ID, Title: "To Delete", Content: "Content"}
	config.DB.Create(&note)

	r := gin.Default()
	routes.SetupNoteRoutes(r)

	// Case 1: Success delete
	req, _ := http.NewRequest(http.MethodDelete, "/api/notes/"+strconv.Itoa(int(note.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var dbNote models.Note
	err := config.DB.First(&dbNote, note.ID).Error
	if err == nil {
		t.Errorf("Note was not deleted from database")
	}

	// Case 2: Not found
	reqNotFound, _ := http.NewRequest(http.MethodDelete, "/api/notes/999", nil)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", wNotFound.Code)
	}
}
