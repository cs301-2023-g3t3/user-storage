package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-storage/controllers"
	"user-storage/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var healthController = new(controllers.HealthController)
var userController = new(controllers.UserController)

type User struct {
    Id        string    `json:"id"`
    FirstName string    `json:"firstName" validate:"required"`
    LastName  string    `json:"lastName" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    Role      *uint       `json:"role" gorm:"default:null"`
}

func SetUpRouter() *gin.Engine{
    router := gin.Default()
    return router
}

func TestHealthCheck(t *testing.T) {
    mockRes :=  "Success"
    r := SetUpRouter()

    r.GET("/health", healthController.CheckHealth)
    req, _ := http.NewRequest("GET", "/health", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    responseData, _ := ioutil.ReadAll(w.Body)
    assert.Equal(t, mockRes, string(responseData))
    assert.Equal(t, http.StatusOK, w.Code)
}

func TestInsertOne(t *testing.T) {
    mockDB, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Failed to create mock DB: %v", err)
    }
    defer mockDB.Close()

    // Open the mock database using gorm.
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn: mockDB,
        SkipInitializeWithVersion: true,
    }), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to open mock database: %v", err)
    }

    userController := controllers.NewUserController(*gormDB)

    expectedQuery := "INSERT INTO users"
    mock.ExpectExec(expectedQuery).
        WithArgs("John", "Doe", "johndoe@gmail.com", 2).
        WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate a successful insert

    userPayload := models.User{
        FirstName: "John",
        LastName:  "Doe",
        Email:     "johndoe@gmail.com",
    }

    // Create a mock gin.Context for testing
    ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

    // Serialize your userPayload to JSON
    userPayloadJSON, err := json.Marshal(userPayload)
    if err != nil {
        t.Fatalf("Failed to marshal userPayload: %v", err)
    }

    // Create a mock request with the userPayload JSON as the request body
    req, _ := http.NewRequest("POST", "/users/accounts", bytes.NewReader(userPayloadJSON))
    req.Header.Set("Content-Type", "application/json")
    ctx.Request = req

    // Call the function being tested
    userController.AddUser(ctx)
    // if err != nil {
    //     t.Errorf("Failed to add user: %v", err)
    // }
    //
    // // Assert that the expected SQL query was executed
    // if err := mock.ExpectationsWereMet(); err != nil {
    //     t.Errorf("SQL expectations were not met: %v", err)
    // }
}
