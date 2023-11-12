package services

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"
	"user-storage/models"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var columns = []string{"id", "first_name", "last_name", "email", "role"}

var gormDB, mock = SetUpDB()

func SetUpDB() (*gorm.DB, sqlmock.Sqlmock){
    // Create a new GORM DB instance with a mocked SQL database
    db, mock, err := sqlmock.New()
    if err != nil {
        log.Fatalf("Error creating mock DB: %v", err)
    }


    mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7.34"))
    // Create a GORM DB connection with the MySQL driver
    gormDB, err := gorm.Open(mysql.New(mysql.Config{
        Conn:                      db,
		DriverName:                "mysql",
		SkipInitializeWithVersion: false,
    }), &gorm.Config{
            SkipDefaultTransaction: true,
        })
    if err != nil {
        log.Fatalf("Error creating GORM DB: %v", err)
    }

    gormDB.AutoMigrate(&models.User{})

    // Insert multiple mock user data into the database
    for i := 1; i <= 10; i++ {
        gormDB.Create(&models.User{
            Id: fmt.Sprint(i),
            FirstName: fmt.Sprintf("John%d", i),
            LastName: "Doe",
            Email: fmt.Sprintf("john%d@example.com", i),
        })
    }
    
    return gormDB, mock
}

// func TestGetAllUsers(t *testing.T) {

//     userService := NewUserService(gormDB)

//     rows := sqlmock.NewRows(columns)

//     var expectedUsers []models.User

//     for i := 1; i <= 10; i++ {
//         role := uint(i%4)
//         user := models.User{
//             Id: fmt.Sprint(i),
//             FirstName: fmt.Sprintf("John%d", i),
//             LastName: "Doe",
//             Email: fmt.Sprintf("john%d@example.com", i),
//             Role: &role,
//         }

//         expectedUsers = append(expectedUsers, user)
//     }

//     for _, user := range expectedUsers {
//         rows.AddRow(user.Id, user.FirstName, user.LastName, user.Email, *user.Role)
//     }

//     statement := "SELECT * FROM `users`"
//     mock.ExpectQuery(regexp.QuoteMeta(statement)).WillReturnRows(rows)

//     users, statusCode, err := userService.GetAllUsers()

//     assert.NoError(t, err) 
//     assert.Equal(t, http.StatusOK, statusCode)
//     assert.Equal(t, expectedUsers, *users)
// }

func TestGetUserById(t *testing.T) {
    
    userService := NewUserService(gormDB)

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
    
    userService := NewUserService(gormDB)
    invalidId := "6969"

    statement := "SELECT * FROM `users` WHERE id = ?"

    mock.ExpectQuery(regexp.QuoteMeta(statement)).
        WithArgs(invalidId).
        WillReturnError(gorm.ErrRecordNotFound)

    res, statusCode, err := userService.GetUserByID(invalidId)
    fmt.Println(err.Error())

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

    userService := NewUserService(gormDB)
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

    userService := NewUserService(gormDB)

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


    mock.ExpectBegin()
    statement := "SELECT * FROM `users` WHERE id = ?"
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

    userService := NewUserService(gormDB)
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

    userService := NewUserService(gormDB)

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

    userService := NewUserService(gormDB)

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
