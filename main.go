package main

import (
	"fmt"
	"log"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Koneksi ke MongoDB
	fmt.Println("🔄 Connecting to MongoDB...")
	config.ConnectDB()
	fmt.Println("✅ MongoDB connected")

	// 2. Setup router
	r := gin.Default()
	routes.SetupRoutes(r)

	// 3. Tambahin endpoint cek server
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// 4. Run server
	port := ":8080"
	fmt.Println("🚀 Server running at http://localhost" + port)
	if err := r.Run(port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
