package controllers

import (
	"net/http"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/application"
	"github.com/gin-gonic/gin"
)

type UpdateProgressController struct {
	useCase *app.UpdateProgress
}

func NewUpdateProgressController(useCase *app.UpdateProgress) *UpdateProgressController {
	return &UpdateProgressController{useCase: useCase}
}

type updateProgressRequest struct {
	Progress int `json:"progress" binding:"required,min=0"`
}

func (ctrl *UpdateProgressController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	retoIDParam := c.Param("retoId")
	retoID, err := strconv.ParseInt(retoIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reto inválido"})
		return
	}

	var req updateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	result, err := ctrl.useCase.Execute(userID.(int64), retoID, req.Progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message":   "Progreso actualizado exitosamente",
		"completed": result.Completed,
	}

	if len(result.LogrosAwarded) > 0 {
		response["logros_awarded"] = result.LogrosAwarded
	}

	c.JSON(http.StatusOK, response)
}
