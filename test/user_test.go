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
	"gorm.io/gorm"
)

var columns = []string{"id", "first_name", "last_name", "email", "role"}
var gormDB, mock = SetUpDB()

func TestGetAllUsers(t *testing.T) {

    userService := services.NewUserService(gormDB)

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

    statement := "SELECT * FROM `users`"
    mock.ExpectQuery(regexp.QuoteMeta(statement)).WillReturnRows(rows)

    users, statusCode, err := userService.GetAllUsers()

    assert.NoError(t, err) 
    assert.Equal(t, http.StatusOK, statusCode)
    assert.Equal(t, expectedUsers, *users)
}

func TestGetUserById(t *testing.T) {
    
    userService := services.NewUserService(gormDB)

    role := uint(1)
    expectedUser := models.User{
        Id: "1",
        FirstName: "John1",
        LastName: "Doe",
        Email: "john1@example.com",
        Role: &role,
    }
    row := sqlmock.NewRows(columns).AddRow("1", "John1", "Doe", "john1@example.com", 1)

    statement := "SELECT * FROM `users` WHERE id = ?"

    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs("1").
        WillReturnRows(row)

    res, statusCode, err := userService.GetUserByID(expectedUser.Id)

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, statusCode) 
    assert.Equal(t, expectedUser, *res)
}

func TestGetUserById_NotFound(t *testing.T) {
    
    userService := services.NewUserService(gormDB)
    invalidId := "6969"

    statement := "SELECT * FROM `users` WHERE id = ?"

    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs(invalidId).
        WillReturnError(gorm.ErrRecordNotFound)

    res, statusCode, err := userService.GetUserByID(invalidId)

    assert.Error(t, err, gorm.ErrRecordNotFound)
    assert.Equal(t, http.StatusNotFound, statusCode) 
    assert.Nil(t, res)
}

func TestAddUser_Success(t *testing.T) {

    firstName, lastName, email, role := "Marilyn", "Monroe", "marilyn@monroe.com", uint(2)

    statement := "INSERT INTO `users` (`id`,`first_name`,`last_name`,`email`,`role`) VALUES (?,?,?,?,?)"

    mock.ExpectBegin()
    mock.ExpectExec(regexp.QuoteMeta(statement)).
		WithArgs(sqlmock.AnyArg(), firstName, lastName, email, role).
		WillReturnResult(sqlmock.NewResult(1, 0))
    mock.ExpectCommit()

    user := models.User{
        FirstName: firstName,
        LastName:  lastName,
        Email:     email,
        Role:      &role,
    }

    userService := services.NewUserService(gormDB)
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

func TestAddUser_BadRequest(t *testing.T) {
    firstName, lastName, role := "Marilyn", "Monroe", uint(2)

    userService := services.NewUserService(gormDB)

    user := models.User{
        FirstName: firstName,
        LastName:  lastName,
        Role:      &role,
    }

    res, statusCode, err := userService.AddUser(&user)

    assert.Error(t, err)
    assert.Equal(t, http.StatusBadRequest, statusCode)
    assert.Nil(t, res)
}

func TestUpdateUserById(t *testing.T) {

    firstName, lastName, email, role := "Marilyn", "Monroe", "marilyn@monroe.com", uint(2)

    row := sqlmock.NewRows(columns).AddRow("1", "John1", "Doe", "john1@example.com", 1)


    statement := "SELECT * FROM `users` WHERE id = ?"
    mock.ExpectBegin()
    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs("1").
        WillReturnRows(row)

    statement = "UPDATE `users` SET `id`=?,`first_name`=?,`last_name`=?,`email`=?,`role`=? WHERE `id` = ?"
    id := "1"

    mock.ExpectExec(regexp.QuoteMeta(statement)).
		WithArgs(id, firstName, lastName, email, role, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
    mock.ExpectCommit()

    user := models.User{
        FirstName: firstName,
        LastName:  lastName,
        Email:     email,
        Role:      &role,
    }

    userService := services.NewUserService(gormDB)
    res, statusCode, err := userService.UpdateUserById(&user, id)
    if err != nil {
        t.Fatal(err)
    }

    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, statusCode)
    assert.Equal(t, firstName, res.FirstName)
    assert.Equal(t, lastName, res.LastName)
    assert.Equal(t, email, res.Email)
    assert.Equal(t, role, *res.Role)
}

func TestUpdateUserById_BadRequest(t *testing.T) {
    firstName, lastName, role := "Marilyn", "Monroe", uint(2)

    userService := services.NewUserService(gormDB)

    user := models.User{
        FirstName: firstName,
        LastName:  lastName,
        Role:      &role,
    }

    res, statusCode, err := userService.UpdateUserById(&user, "1")

    assert.Error(t, err)
    assert.Equal(t, http.StatusBadRequest, statusCode)
    assert.Nil(t, res)
}

func TestUpdateUserById_NotFound(t *testing.T) {
    firstName, lastName, email, role := "Marilyn", "Monroe", "john1@example.com", uint(2)
    invalidId := "6969"

    statement := "SELECT * FROM `users` WHERE id = ?"

    mock.ExpectBegin()
    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs(invalidId).
        WillReturnError(gorm.ErrRecordNotFound)
    mock.ExpectCommit()

    userService := services.NewUserService(gormDB)

    user := models.User{
        FirstName: firstName,
        LastName:  lastName,
        Email:     email,
        Role:      &role,
    }

    res, statusCode, err := userService.UpdateUserById(&user, invalidId)
    fmt.Println(err.Error())

    assert.Error(t, err, gorm.ErrRecordNotFound)
    assert.Equal(t, http.StatusNotFound, statusCode)
    assert.Nil(t, res)
}

// func TestDeleteUserById(t *testing.T){
//     id, firstName, lastName, email, role := "1", "John1", "Doe", "john1@example.com", uint(1)
//     row := sqlmock.NewRows(columns).AddRow(id, firstName, lastName, email, role)
//
//     mock.ExpectBegin()
//     mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ?")).
//         WithArgs(id).
//         WillReturnRows(row)
//
//     mock.ExpectQuery(regexp.QuoteMeta("DELETE FROM `users`")).
//         WithArgs(id).
//         WillReturnRows(row)
//     mock.ExpectCommit()
//
//     userService := services.NewUserService(gormDB)
//     res, statusCode, err := userService.DeleteUserById(id)
//     if err != nil {
//         t.Fatal(err)
//     }
//
//     assert.NoError(t, err)
//     assert.Equal(t, http.StatusOK, statusCode)
//     assert.Equal(t, "Success", *res)
// }
