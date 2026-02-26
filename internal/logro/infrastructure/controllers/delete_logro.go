package controllers

import (
	"net/http"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/logro/application"
	"github.com/gin-gonic/gin"
)

type DeleteLogroController struct {
	useCase *app.DeleteLogro
}

func NewDeleteLogroController(useCase *app.DeleteLogro) *DeleteLogroController {
	return &DeleteLogroController{useCase: useCase}
}

func (ctrl *DeleteLogroController) Handle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = ctrl.useCase.Execute(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logro eliminado exitosamente",
	})
}
