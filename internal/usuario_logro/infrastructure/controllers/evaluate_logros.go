package controllers

import (
	"net/http"

	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/application"
	"github.com/gin-gonic/gin"
)

type EvaluateLogrosController struct {
	useCase *app.EvaluateLogros
}

func NewEvaluateLogrosController(useCase *app.EvaluateLogros) *EvaluateLogrosController {
	return &EvaluateLogrosController{useCase: useCase}
}

// Handle evalúa y otorga automáticamente los logros que el usuario
// autenticado haya cumplido. Se invoca internamente después de completar
// un reto o al actualizar puntos.
func (ctrl *EvaluateLogrosController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	awarded, err := ctrl.useCase.Execute(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Evaluación de logros completada",
		"logros_awarded": awarded,
	})
}
