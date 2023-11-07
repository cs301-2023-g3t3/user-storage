package main

import (
	"net/http"
	"testing"
	"user-storage/models"
	"user-storage/services"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
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

    return gormDB, mock
}

func TestGetAllUsers(t *testing.T) {
    gormDB, mock := SetUpDB(t)

    userService := services.NewUserService(gormDB)

    // Define your expected data
    role := uint(4)
    expectedUsers := []models.User{
        {Id: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com", Role: &role},
        {Id: "2", FirstName: "Yoa", LastName: "Lung", Email: "yoa@lung.com", Role: &role},
    }

    // Define your expected rows and columns
    columns := []string{"id", "first_name", "last_name", "email", "role"}
    expectedRows := sqlmock.NewRows(columns).
        AddRow("1", "John", "Doe", "john@example.com", 4).
        AddRow("2", "Yoa", "Lung", "yoa@lung.com", 4)

    // Set up the mock expectations for the SQL query
    mock.ExpectQuery("SELECT * FROM `users`").WillReturnRows(expectedRows)

    // Call the GetAllUsers method
    users, statusCode, err := userService.GetAllUsers()

    // Assertions
    assert.NoError(t, err) // Ensure no error occurred
    assert.Equal(t, http.StatusOK, statusCode) // Check the HTTP status code
    assert.Equal(t, expectedUsers, *users) // Compare the expected users with the returned users
}
