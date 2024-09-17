package controllers

import (
	"arithmetic-calculator/config"
	"arithmetic-calculator/controllers"
	"arithmetic-calculator/models"
	"arithmetic-calculator/routes"
	"arithmetic-calculator/services"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter() *gin.Engine {
	config.ConnectDatabase()
	return routes.SetupRouter()
}

func TestLoginSuccess(t *testing.T) {
	r := SetupTestRouter()

	// Assume user already exists in DB

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"username":"testuser","password":"password"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "token")
}

func TestLoginFailure(t *testing.T) {
	r := SetupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"username":"wronguser","password":"wrongpassword"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestFetchRecordsWithPagination(t *testing.T) {
	config.ConnectDatabase()

	page := 1
	pageSize := 5

	paginatedRecords, err := services.FetchRecords("testuser", page, pageSize)

	assert.Nil(t, err)
	assert.Equal(t, page, paginatedRecords.Page)
	assert.NotNil(t, paginatedRecords.Records)
	assert.True(t, paginatedRecords.TotalItems > 0)
	assert.True(t, paginatedRecords.TotalPages >= 1)
}

func TestRegister(t *testing.T) {
	config.ConnectDatabase()

	router := gin.Default()
	router.POST("/api/v1/auth/register", controllers.Register)

	user := models.User{
		Username: "newuser@example.com",
		Password: "password123",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")
}
