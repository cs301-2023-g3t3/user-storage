package main

import (
	"user-storage/cache"
	"user-storage/models"

	log "github.com/sirupsen/logrus"
)

func init() {
	models.ConnectToDB()
	cache.ConnectToRedis()
	log.SetFormatter(&log.JSONFormatter{})

}

func main() {
	InitRoutes()
}
