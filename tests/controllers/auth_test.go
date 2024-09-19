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
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupTestRouter() *gin.Engine {
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_USER", "username")
	os.Setenv("DB_PASSWORD", "password")
	os.Setenv("DB_NAME", "arithmetic_db")
	os.Setenv("DB_PORT", "5432")

	config.InitDB()

	return routes.SetupRouter()
}

func TestLoginSuccess(t *testing.T) {
	r := SetupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(`{"username":"testuser1@example.com","password":"password123"}`))
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
	assert.Contains(t, w.Body.String(), "user not found")
}

func TestFetchRecordsWithPagination(t *testing.T) {
	SetupTestRouter()

	page := 1
	pageSize := 5

	paginatedRecords, err := services.FetchRecords("testuser1@example.com", page, pageSize)

	assert.Nil(t, err)
	assert.Equal(t, page, paginatedRecords.Page)
	assert.NotNil(t, paginatedRecords.Records)
	assert.True(t, paginatedRecords.TotalItems > 0)
	assert.True(t, paginatedRecords.TotalPages >= 1)
}

func TestRegister(t *testing.T) {
	_ = SetupTestRouter()

	router := gin.Default()
	router.POST("/api/v1/auth/register", controllers.Register)

	existingUser := models.User{}
	config.DB.Unscoped().Where("username = ?", "newuser123@example.com").Delete(&existingUser)

	user := models.User{
		Username: "newuser123@example.com",
		Password: "password123",
	}

	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "User registered successfully")

	var registeredUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&registeredUser).Error; err != nil {
		t.Fatalf("Failed to find the user: %v", err)
	}

	if err := config.DB.Delete(&registeredUser).Error; err != nil {
		t.Fatalf("Failed to delete the user: %v", err)
	}
}

func TestRegisterExistingUser(t *testing.T) {
	r := SetupTestRouter()

	existingUser := models.User{
		Username: "testuser1@example.com",
		Password: "password123",
	}

	jsonValue, _ := json.Marshal(existingUser)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "User already exists")
}

func TestGetProfile(t *testing.T) {
	_ = SetupTestRouter()

	router := gin.Default()
	router.POST("/api/v1/auth/login", controllers.Login)

	loginPayload := `{"username":"testuser1@example.com","password":"password123"}`
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", strings.NewReader(loginPayload))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var loginResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	if err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	_, ok := loginResponse["token"]
	if !ok {
		t.Fatalf("Token not found in login response")
	}

	assert.Contains(t, loginResponse["username"], "testuser1@example.com")
}
