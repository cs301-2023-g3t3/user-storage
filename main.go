package main

import (
	log "github.com/sirupsen/logrus"
	"user-storage/models"
)

func init() {
	models.ConnectToDB()
	log.SetFormatter(&log.JSONFormatter{})

}

func main() {
	InitRoutes()
}
