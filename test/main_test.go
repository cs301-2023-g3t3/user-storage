package main

import (
	"fmt"
	"log"
	"user-storage/models"

    sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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
