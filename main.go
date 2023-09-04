package main

import "user-storage/models"

func init() {
	models.ConnectToDB()
}

func main() {
	InitRoutes()
}
