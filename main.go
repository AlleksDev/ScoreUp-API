package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Wire inyecta TODAS las dependencias automáticamente.
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Error al inicializar la aplicación: %v", err)
	}

	// Iniciar el Hub de WebSocket en una goroutine separada
	go app.Hub.Run()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Crear el servidor HTTP manualmente para poder hacer shutdown graceful.
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: app.Engine,
	}

	// Arrancar el servidor en una goroutine.
	go func() {
		log.Printf("Servidor ScoreUp-API iniciado en http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v", err)
		}
	}()

	// Esperar señal de apagado (Ctrl+C, docker stop, etc.)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Apagando servidor...")

	// Dar hasta 10 segundos para que las peticiones en curso terminen.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Apagar Hub de WebSocket (cierra todas las conexiones WS).
	app.Hub.Shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error durante shutdown del servidor: %v", err)
	}

	log.Println("Servidor detenido correctamente")
}
