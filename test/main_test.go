package main

import (
	"fmt"
	"io"
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

func SetUpDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock){
    // Create a new GORM DB instance with a mocked SQL database
    db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
    if err != nil {
        t.Fatalf("Error creating mock DB: %v", err)
    }

    // Create a GORM DB connection with the MySQL driver
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn: db,
        SkipInitializeWithVersion: true,
    }), &gorm.Config{})
    if err != nil {
        t.Fatalf("Error creating GORM DB: %v", err)
    }

    gormDB.AutoMigrate(&models.User{})
    // Insert the mock data
    roleVal := uint(1)
    gormDB.Create(&models.User{
        Id: "1",
        FirstName: "John",
        LastName: "Doe",
        Email: "john@example.com",
        Role: &roleVal,
    })

    return gormDB, mock
}

func TestGetAllUsers(t *testing.T) {
    gormDB, mock := SetUpDB(t)

    userController := controllers.NewUserController(*gormDB)

    req, err := http.NewRequest("GET", "/users/accounts", nil)
    if err != nil {
        t.Fatalf("an error '%s' was not expected while creating request", err)
    }

    recorder := httptest.NewRecorder()

    ctx, _ := gin.CreateTestContext(recorder)
    ctx.Request = req
    ctx.Params = gin.Params{}
    columns := []string{"id", "first_name", "last_name", "email", "role"}
    expectedRows := sqlmock.NewRows(columns).AddRow("1", "John", "Doe", "john@example.com", 4)

    // Define your expected rows and columns
    // Set up the mock expectations
    mock.ExpectQuery("SELECT * FROM `users`").WillReturnRows(expectedRows)

    userController.GetAllUsers(ctx)

    responseBody, err := io.ReadAll(recorder.Body)
    if err != nil {
        t.Fatalf("error reading respons body: %v", err)
    }

    responseBodyStr := string(responseBody)
    fmt.Println(responseBodyStr)

    // Assertions
    assert.Equal(t, http.StatusOK, ctx.Writer.Status()) // Check the HTTP status code
}
