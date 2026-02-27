package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	logroInfra "github.com/AlleksDev/ScoreUp-API/internal/logro/infrastructure"
	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
	retoInfra "github.com/AlleksDev/ScoreUp-API/internal/reto/infrastructure"
	userInfra "github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure"
	usuarioLogroInfra "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/infrastructure"
	usuarioRetoInfra "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/infrastructure"
	"github.com/AlleksDev/ScoreUp-API/internal/websocket"
)

// App es el contenedor raíz de la aplicación. Wire la construye automáticamente.
type App struct {
	Engine *gin.Engine
	Hub    *websocket.Hub
}

// ProvideGinEngine es un provider de Wire que recibe TODOS los módulos ya
// inyectados y registra sus rutas.
func ProvideGinEngine(
	userMod *userInfra.UserModule,
	retoMod *retoInfra.RetoModule,
	logroMod *logroInfra.LogroModule,
	usuarioRetoMod *usuarioRetoInfra.UsuarioRetoModule,
	usuarioLogroMod *usuarioLogroInfra.UsuarioLogroModule,
	wsHandler *websocket.WSHandler,
) *gin.Engine {
	r := gin.Default()
	r.Use(core.SetupCORS())

	// Rutas públicas (sin autenticación)
	userMod.RegisterRoutes(r)

	// WebSocket endpoint
	r.GET("/ws", wsHandler.HandleConnection)

	// Rutas protegidas (con JWT auth middleware)
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		retoMod.RegisterRoutes(api)
		logroMod.RegisterRoutes(api)
		usuarioRetoMod.RegisterRoutes(api)
		usuarioLogroMod.RegisterRoutes(api)
	}

	return r
}

// NewApp es el provider final que Wire usa para ensamblar la aplicación.
func NewApp(engine *gin.Engine, hub *websocket.Hub) *App {
	return &App{Engine: engine, Hub: hub}
}
