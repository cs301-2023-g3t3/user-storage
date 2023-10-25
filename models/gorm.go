// From points-ledger/models/gorm.go
package models

import (
	"fmt"
	"log"
	"os"

    "github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	env := os.Getenv("ENV")
	if env != "lambda" {
        err = godotenv.Load()
        switch os.Getenv("DB_TYPE") {
            case "local":
                err = godotenv.Load(".env.local")
            case "live":
                err = godotenv.Load(".env.live")
            }
            if err != nil {
                log.Fatal("Error loading .env file")
            }
    }
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, hostname, port, dbname) + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to Database")
	}
}
