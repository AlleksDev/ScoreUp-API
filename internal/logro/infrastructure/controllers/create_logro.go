package controllers

import (
	"net/http"

	app "github.com/AlleksDev/ScoreUp-API/internal/logro/application"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/domain/entities"
	"github.com/gin-gonic/gin"
)

type CreateLogroController struct {
	useCase *app.CreateLogro
}

func NewCreateLogroController(useCase *app.CreateLogro) *CreateLogroController {
	return &CreateLogroController{useCase: useCase}
}

type createLogroRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description" binding:"required"`
	RequiredPoints int    `json:"required_points"`
	RequiredRetos  int    `json:"required_retos"`
}

func (ctrl *CreateLogroController) Handle(c *gin.Context) {
	var req createLogroRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	logro := entities.Logro{
		Name:           req.Name,
		Description:    req.Description,
		RequiredPoints: req.RequiredPoints,
		RequiredRetos:  req.RequiredRetos,
	}

	err := ctrl.useCase.Execute(&logro)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Logro creado exitosamente",
		"id":      logro.ID,
	})
}
