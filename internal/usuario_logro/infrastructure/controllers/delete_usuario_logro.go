package controllers

import (
	"net/http"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/application"
	"github.com/gin-gonic/gin"
)

type DeleteUsuarioLogroController struct {
	useCase *app.DeleteUsuarioLogro
}

func NewDeleteUsuarioLogroController(useCase *app.DeleteUsuarioLogro) *DeleteUsuarioLogroController {
	return &DeleteUsuarioLogroController{useCase: useCase}
}

func (ctrl *DeleteUsuarioLogroController) Handle(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	logroIDParam := c.Param("logroId")
	logroID, err := strconv.ParseInt(logroIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de logro inválido"})
		return
	}

	err = ctrl.useCase.Execute(userID.(int64), logroID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logro removido del usuario exitosamente",
	})
}
