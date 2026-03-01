package controllers

import (
	"net/http"
	"strconv"

	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/application"
	"github.com/gin-gonic/gin"
)

type GetUsuarioRetosController struct {
	useCase *app.GetUsuarioRetos
}

func NewGetUsuarioRetosController(useCase *app.GetUsuarioRetos) *GetUsuarioRetosController {
	return &GetUsuarioRetosController{useCase: useCase}
}

func (ctrl *GetUsuarioRetosController) HandleByUser(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	results, err := ctrl.useCase.ExecuteByUser(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuario_retos": results,
	})
}

func (ctrl *GetUsuarioRetosController) HandleByReto(c *gin.Context) {
	retoIDParam := c.Param("retoId")
	retoID, err := strconv.ParseInt(retoIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reto inválido"})
		return
	}

	results, err := ctrl.useCase.ExecuteByReto(retoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"usuario_retos": results,
	})
}
