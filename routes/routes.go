package routes

import (
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Auth routes
	r.GET("/auth/google/login", controllers.Login)
	r.GET("/auth/google/callback", controllers.GoogleCallback)

	// Contoh proteksi route pakai middleware JWT nanti bisa ditambah
	r.GET("/profile", controllers.ProfileHandler)

	// Register user
	r.POST("/register", controllers.Register)
}
