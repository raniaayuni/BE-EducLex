package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Route GET /
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Halo dari EducLex 🚀",
		})
	})

	// Jalankan server di port 8080
	r.Run(":8080")
}
