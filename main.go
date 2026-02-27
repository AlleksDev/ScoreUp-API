package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Wire inyecta TODAS las dependencias automáticamente.
	// InitializeApp() está generado en wire_gen.go
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Error al inicializar la aplicación: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor ScoreUp-API iniciado en http://localhost:%s", port)
	if err := app.Engine.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
