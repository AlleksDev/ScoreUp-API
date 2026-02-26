package controllers

import (
	"net/http"

	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/application"
	"github.com/gin-gonic/gin"
)

type GetUsuarioLogrosController struct {
	useCase *app.GetUsuarioLogros
}

func NewGetUsuarioLogrosController(useCase *app.GetUsuarioLogros) *GetUsuarioLogrosController {
	return &GetUsuarioLogrosController{useCase: useCase}
}

func (ctrl *GetUsuarioLogrosController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	logros, err := ctrl.useCase.Execute(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuario_logros": logros,
	})
}
