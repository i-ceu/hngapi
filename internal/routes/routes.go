package routes

import (
	"profile-api/internal/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() {

	r := gin.Default()
	{
		r.GET("/me", controllers.GetProfile)
		r.POST("/strings", controllers.AddStrings)
		r.GET("/strings/:string_value", controllers.GetString)
		r.GET("/strings", controllers.GetAllStrings)
		r.GET("/strings/filter-by-natural-language", controllers.FilterByNaturalLanguage)
		r.DELETE("/strings:string_value", controllers.DeleteString)
	}

	r.Run()
}
