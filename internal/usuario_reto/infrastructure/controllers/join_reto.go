package controllers

import (
	"net/http"
	"strconv"

	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/application"
	"github.com/gin-gonic/gin"
)

type JoinRetoController struct {
	useCase *app.JoinReto
}

func NewJoinRetoController(useCase *app.JoinReto) *JoinRetoController {
	return &JoinRetoController{useCase: useCase}
}

type joinRetoRequest struct {
	RetoID int64 `json:"reto_id" binding:"required"`
}

func (ctrl *JoinRetoController) Handle(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	var req joinRetoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	err := ctrl.useCase.Execute(uid, req.RetoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Unido al reto exitosamente",
	})
}

// HandleByParam une al usuario al reto indicado por parámetro de ruta
func (ctrl *JoinRetoController) HandleByParam(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	retoIDParam := c.Param("retoId")
	retoID, err := strconv.ParseInt(retoIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reto inválido"})
		return
	}

	err = ctrl.useCase.Execute(uid, retoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Unido al reto exitosamente",
	})
}
