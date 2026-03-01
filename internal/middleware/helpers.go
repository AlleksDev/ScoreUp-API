package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserID extrae de forma segura el userID del contexto de Gin.
// Devuelve el ID y true si existe y es del tipo correcto; de lo contrario
// responde automáticamente con 401 y devuelve (0, false).
func GetUserID(c *gin.Context) (int64, bool) {
	val, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return 0, false
	}
	uid, ok := val.(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Tipo de userID inválido en contexto"})
		return 0, false
	}
	return uid, true
}
