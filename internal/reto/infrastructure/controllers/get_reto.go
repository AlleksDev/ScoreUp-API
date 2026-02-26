package controllers

import (
	"net/http"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/reto/application"
	"github.com/gin-gonic/gin"
)

type GetRetoController struct {
	useCase *app.GetReto
}

func NewGetRetoController(useCase *app.GetReto) *GetRetoController {
	return &GetRetoController{useCase: useCase}
}

func (ctrl *GetRetoController) HandleGetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	reto, err := ctrl.useCase.ExecuteByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reto": reto,
	})
}

func (ctrl *GetRetoController) HandleGetByUserID(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	retos, err := ctrl.useCase.ExecuteByUserID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"retos": retos,
	})
}
