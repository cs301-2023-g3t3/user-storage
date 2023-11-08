package main

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"user-storage/models"
	"user-storage/services"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {
    gormDB, _, mock := SetUpDB()

    userService := services.NewUserService(gormDB)

    columns := []string{"id", "first_name", "last_name", "email", "role"}
    rows := sqlmock.NewRows(columns)

    var expectedUsers []models.User

    for i := 1; i <= 10; i++ {
        role := uint(i%4)
        user := models.User{
            Id: fmt.Sprint(i),
            FirstName: fmt.Sprintf("John%d", i),
            LastName: "Doe",
            Email: fmt.Sprintf("john%d@example.com", i),
            Role: &role,
        }

        expectedUsers = append(expectedUsers, user)
    }

    for _, user := range expectedUsers {
        rows.AddRow(user.Id, user.FirstName, user.LastName, user.Email, *user.Role)
    }

    mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).WillReturnRows(rows)

    users, statusCode, err := userService.GetAllUsers()

    assert.NoError(t, err) 
    assert.Equal(t, http.StatusOK, statusCode)
    assert.Equal(t, expectedUsers, *users)
}

func TestAddUser_Success(t *testing.T) {
    gormDB, db, mock := SetUpDB()
    db.Begin()

    firstName, lastName, email, role := "Marilyn", "Monroe", "marilyn@monroe.com", uint(2)

    userService := services.NewUserService(gormDB)

    mock.ExpectBegin()
    mock.ExpectExec("INSERT INTO `users`").
		WithArgs(sqlmock.AnyArg(), firstName, lastName, email, role).
		WillReturnResult(sqlmock.NewResult(1, 0))
    mock.ExpectCommit()

    user := models.User{
        FirstName: firstName,
        LastName:  lastName,
        Email:     email,
        Role:      &role,
    }

    res, statusCode, err := userService.AddUser(&user)
    if err != nil {
        t.Fatalf("Error getting user by ID: %v", err)
    }

    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, statusCode)
    assert.Equal(t, firstName, res.FirstName)
    assert.Equal(t, lastName, res.LastName)
    assert.Equal(t, email, res.Email)
    assert.Equal(t, role, *res.Role)
}
