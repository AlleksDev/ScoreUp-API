package controllers

import (
	"net/http"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/logro/application"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"
	"github.com/gin-gonic/gin"
)

type UpdateLogroController struct {
	useCase *app.UpdateLogro
}

func NewUpdateLogroController(useCase *app.UpdateLogro) *UpdateLogroController {
	return &UpdateLogroController{useCase: useCase}
}

type updateLogroRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	RequiredPoints int    `json:"required_points"`
	RequiredRetos  int    `json:"required_retos"`
}

func (ctrl *UpdateLogroController) Handle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req updateLogroRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	logro := entities.Logro{
		ID:             id,
		Name:           req.Name,
		Description:    req.Description,
		RequiredPoints: req.RequiredPoints,
		RequiredRetos:  req.RequiredRetos,
	}

	err = ctrl.useCase.Execute(&logro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logro actualizado exitosamente",
	})
}
