package core

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupCORS configura CORS de forma permisiva.
// Las apps móviles nativas NO están sujetas a CORS (es un mecanismo del navegador),
// pero esta configuración permite también consumir la API desde webapps y herramientas de testing.
func SetupCORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			// Aceptar cualquier origen. En producción se puede restringir
			// a dominios concretos si se sirve también una webapp.
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "Accept", "Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
