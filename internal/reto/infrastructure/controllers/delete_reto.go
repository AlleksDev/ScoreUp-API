package controllers

import (
	"net/http"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/reto/application"
	"github.com/gin-gonic/gin"
)

type DeleteRetoController struct {
	useCase *app.DeleteReto
}

func NewDeleteRetoController(useCase *app.DeleteReto) *DeleteRetoController {
	return &DeleteRetoController{useCase: useCase}
}

func (ctrl *DeleteRetoController) Handle(c *gin.Context) {
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
		"message": "Reto eliminado exitosamente",
	})
}
