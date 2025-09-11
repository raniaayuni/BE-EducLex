package main

import (
	"log"

	"github.com/EducLex/BE-EducLex/config"
	"github.com/EducLex/BE-EducLex/routes"
)

func main() {
	// koneksi DB
	config.ConnectDB()

	// setup router
	r := routes.SetupRouter()
	log.Println("Server running on :8080")

	// run server
	r.Run(":8080")
}
