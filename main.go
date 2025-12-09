package main

import (
	"log"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/controllers"
	"github.com/EducLex/BE-EducLex/routes"
	"github.com/gin-contrib/cors"
)

func main() {
	// Koneksi database
	config.ConnectDB()

	// Setup router (CORS sudah ada di dalam SetupRouter)
	r := routes.SetupRouter()
	// Konfigurasi CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5501"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Seed kategori
	controllers.SeedCategories()

	// Aktifkan CORS
	r.Use(cors.Default())

	log.Println("Server running on :8080")
	r.Run(":8080")
}
