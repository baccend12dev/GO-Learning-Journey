package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Initialize in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory test database: %v", err)
	}

	// Migrate models
	err = db.AutoMigrate(&models.Server{})
	if err != nil {
		t.Fatalf("Failed to migrate database schemas: %v", err)
	}

	// Override global config.DB with the test database instance
	config.DB = db
}

func TestGetServers(t *testing.T) {
	setupTestDB(t)

	// Seed some test server records
	servers := []models.Server{
		{Name: "Server 1", IP: "192.168.1.10", OS: "Ubuntu 22.04", Location: "Data Center A"},
		{Name: "Server 2", IP: "192.168.1.11", OS: "Rocky Linux 9", Location: "Data Center B"},
	}
	for _, s := range servers {
		config.DB.Create(&s)
	}

	r := gin.Default()
	r.GET("/api/servers", GetServers)

	req, _ := http.NewRequest(http.MethodGet, "/api/servers", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var returnedServers []models.Server
	if err := json.Unmarshal(w.Body.Bytes(), &returnedServers); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(returnedServers) != 2 {
		t.Errorf("Expected 2 servers, but got %d", len(returnedServers))
	}

	if returnedServers[0].Name != "Server 1" || returnedServers[1].Name != "Server 2" {
		t.Errorf("Seeded servers data mismatch")
	}
}

func TestGetServerByID(t *testing.T) {
	setupTestDB(t)

	server := models.Server{Name: "Server 1", IP: "192.168.1.10", OS: "Ubuntu 22.04", Location: "Data Center A"}
	config.DB.Create(&server)

	r := gin.Default()
	r.GET("/api/servers/:id", GetServerByID)

	// Case 1: ID exists
	req, _ := http.NewRequest(http.MethodGet, "/api/servers/"+strconv.Itoa(int(server.ID)), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d for existing server, but got %d", http.StatusOK, w.Code)
	}

	var returnedServer models.Server
	if err := json.Unmarshal(w.Body.Bytes(), &returnedServer); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if returnedServer.ID != server.ID || returnedServer.Name != server.Name {
		t.Errorf("Server details mismatch: expected ID %d and Name %q, got ID %d and Name %q",
			server.ID, server.Name, returnedServer.ID, returnedServer.Name)
	}

	// Case 2: ID does not exist
	reqNotFound, _ := http.NewRequest(http.MethodGet, "/api/servers/999", nil)
	wNotFound := httptest.NewRecorder()
	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for non-existent server, but got %d", http.StatusNotFound, wNotFound.Code)
	}
}

func TestCreateServer(t *testing.T) {
	setupTestDB(t)

	r := gin.Default()
	r.POST("/api/servers", CreateServer)

	// Case 1: Valid input
	reqBody := models.CreateServerRequest{
		Name:     "Test Server",
		IP:       "10.0.0.5",
		OS:       "Debian 12",
		Location: "Rack A1",
	}
	bodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/api/servers", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, but got %d", http.StatusCreated, w.Code)
	}

	var responseServer models.Server
	if err := json.Unmarshal(w.Body.Bytes(), &responseServer); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if responseServer.Name != reqBody.Name || responseServer.IP != reqBody.IP || responseServer.OS != reqBody.OS || responseServer.Location != reqBody.Location {
		t.Errorf("Response server fields mismatch with requested fields")
	}

	// Verify server actually exists in the database
	var dbServer models.Server
	if err := config.DB.First(&dbServer, responseServer.ID).Error; err != nil {
		t.Errorf("Server was not stored in the database: %v", err)
	}

	// Case 2: Invalid input (missing required fields)
	invalidReqBody := map[string]interface{}{
		"name": "Missing Other Fields Server",
	}
	invalidBytes, _ := json.Marshal(invalidReqBody)
	reqInvalid, _ := http.NewRequest(http.MethodPost, "/api/servers", bytes.NewBuffer(invalidBytes))
	reqInvalid.Header.Set("Content-Type", "application/json")
	wInvalid := httptest.NewRecorder()

	r.ServeHTTP(wInvalid, reqInvalid)

	if wInvalid.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d on validation failure, but got %d", http.StatusBadRequest, wInvalid.Code)
	}
}

func TestUpdateServer(t *testing.T) {
	setupTestDB(t)

	server := models.Server{Name: "Old Name", IP: "192.168.1.10", OS: "Ubuntu 20.04", Location: "Data Center A"}
	config.DB.Create(&server)

	r := gin.Default()
	r.PUT("/api/servers/:id", UpdateServer)

	// Case 1: Update success
	updatedServerData := map[string]interface{}{
		"Name":     "New Name",
		"IP":       "192.168.1.20",
		"OS":       "Ubuntu 22.04",
		"Location": "Data Center B",
	}
	bodyBytes, _ := json.Marshal(updatedServerData)
	req, _ := http.NewRequest(http.MethodPut, "/api/servers/"+strconv.Itoa(int(server.ID)), bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	var returnedServer models.Server
	if err := json.Unmarshal(w.Body.Bytes(), &returnedServer); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if returnedServer.Name != updatedServerData["Name"].(string) || returnedServer.IP != updatedServerData["IP"].(string) {
		t.Errorf("Returned server does not match updated values")
	}

	// Verify database changes
	var dbServer models.Server
	config.DB.First(&dbServer, server.ID)
	if dbServer.Name != updatedServerData["Name"].(string) || dbServer.IP != updatedServerData["IP"].(string) {
		t.Errorf("Database server was not updated successfully")
	}

	// Case 2: Update on non-existent server
	reqNotFound, _ := http.NewRequest(http.MethodPut, "/api/servers/999", bytes.NewBuffer(bodyBytes))
	reqNotFound.Header.Set("Content-Type", "application/json")
	wNotFound := httptest.NewRecorder()

	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d for non-existent server update, but got %d", http.StatusNotFound, wNotFound.Code)
	}
}

func TestDeleteServer(t *testing.T) {
	setupTestDB(t)

	server := models.Server{Name: "Server To Delete", IP: "192.168.1.10", OS: "Ubuntu 22.04", Location: "Data Center A"}
	config.DB.Create(&server)

	r := gin.Default()
	r.DELETE("/api/servers/:id", DeleteServer)

	// Case 1: Delete success
	req, _ := http.NewRequest(http.MethodDelete, "/api/servers/"+strconv.Itoa(int(server.ID)), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, w.Code)
	}

	// Verify database shows deleted (First returns record not found error)
	var dbServer models.Server
	err := config.DB.First(&dbServer, server.ID).Error
	if err == nil {
		t.Errorf("Server record was not deleted from the database")
	}

	// Case 2: Delete on non-existent server
	reqNotFound, _ := http.NewRequest(http.MethodDelete, "/api/servers/999", nil)
	wNotFound := httptest.NewRecorder()

	r.ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, but got %d", http.StatusNotFound, wNotFound.Code)
	}
}
