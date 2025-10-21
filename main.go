package main

import (
	"profile-api/internal/config"
	"profile-api/internal/models"
	"profile-api/internal/routes"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToDB()
	config.DB.AutoMigrate(&models.String{})
}

func main() {
	routes.RegisterRoutes()
}
