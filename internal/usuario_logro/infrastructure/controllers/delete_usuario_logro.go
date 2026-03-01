package controllers

import (
	"net/http"
	"strconv"

	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
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
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	logroIDParam := c.Param("logroId")
	logroID, err := strconv.ParseInt(logroIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de logro inválido"})
		return
	}

	err = ctrl.useCase.Execute(uid, logroID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logro removido del usuario exitosamente",
	})
}
