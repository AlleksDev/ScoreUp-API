package controllers

import (
	"net/http"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/logro/application"
	"github.com/gin-gonic/gin"
)

type GetLogroController struct {
	useCase *app.GetLogro
}

func NewGetLogroController(useCase *app.GetLogro) *GetLogroController {
	return &GetLogroController{useCase: useCase}
}

func (ctrl *GetLogroController) HandleGetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	logro, err := ctrl.useCase.ExecuteByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logro": logro,
	})
}

func (ctrl *GetLogroController) HandleGetAll(c *gin.Context) {
	logros, err := ctrl.useCase.ExecuteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logros": logros,
	})
}
