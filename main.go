package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Logging a archivo para capturar errores incluso si stdout se pierde.
	logFile, err := os.OpenFile("scoreup-api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Loguear tanto a stdout como al archivo.
		log.SetOutput(os.Stdout) // se mantiene stdout
		defer logFile.Close()
		// Usar un MultiWriter para duplicar logs.
		log.SetOutput(newMultiWriter(os.Stdout, logFile))
	}

	// Recover de último recurso para panics en el hilo principal.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("[FATAL PANIC] main goroutine crashed: %v\n%s", r, debug.Stack())
		}
	}()

	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Wire inyecta TODAS las dependencias automáticamente.
	app, err := InitializeApp()
	if err != nil {
		log.Fatalf("Error al inicializar la aplicación: %v", err)
	}

	// Iniciar el Hub de WebSocket en una goroutine protegida.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] Hub.Run() crashed: %v\n%s", r, debug.Stack())
			}
		}()
		app.Hub.Run()
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Crear el servidor HTTP manualmente para poder hacer shutdown graceful.
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      app.Engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Canal para señalizar error fatal del servidor.
	serverErr := make(chan error, 1)

	// Arrancar el servidor en una goroutine.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC] HTTP server goroutine crashed: %v\n%s", r, debug.Stack())
			}
		}()
		log.Printf("Servidor ScoreUp-API iniciado en http://localhost:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("[ERROR] Error al iniciar el servidor: %v", err)
			serverErr <- err
		}
	}()

	// Esperar señal de apagado (Ctrl+C, docker stop, etc.) O error fatal del servidor.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("Señal de apagado recibida...")
	case err := <-serverErr:
		log.Printf("Servidor se detuvo por error: %v", err)
	}

	log.Println("Apagando servidor...")

	// Dar hasta 10 segundos para que las peticiones en curso terminen.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Apagar Hub de WebSocket (cierra todas las conexiones WS).
	app.Hub.Shutdown()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Error durante shutdown del servidor: %v", err)
	}

	log.Println("Servidor detenido correctamente")
}

// newMultiWriter crea un io.Writer que escribe en múltiples destinos.
func newMultiWriter(writers ...interface{ Write([]byte) (int, error) }) *multiWriter {
	return &multiWriter{writers: writers}
}

type multiWriter struct {
	writers []interface{ Write([]byte) (int, error) }
}

func (mw *multiWriter) Write(p []byte) (int, error) {
	for _, w := range mw.writers {
		w.Write(p)
	}
	return len(p), nil
}
