package controllers

import (
	"net/http"

	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
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
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	logros, err := ctrl.useCase.Execute(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuario_logros": logros,
	})
}
