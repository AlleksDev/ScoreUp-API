package controllers

import (
	"net/http"
	"strconv"
	"time"

	app "github.com/AlleksDev/ScoreUp-API/internal/reto/application"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"
	"github.com/gin-gonic/gin"
)

type UpdateRetoController struct {
	useCase *app.UpdateReto
}

func NewUpdateRetoController(useCase *app.UpdateReto) *UpdateRetoController {
	return &UpdateRetoController{useCase: useCase}
}

type updateRetoRequest struct {
	Subject       string  `json:"subject" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Goal          int     `json:"goal" binding:"required,min=1"`
	Progress      int     `json:"progress"`
	PointsAwarded int     `json:"points_awarded"`
	Deadline      *string `json:"deadline"`
	Status        string  `json:"status" binding:"required,oneof=activo completado"`
}

func (ctrl *UpdateRetoController) Handle(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req updateRetoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	reto := entities.Reto{
		ID:            id,
		Subject:       req.Subject,
		Description:   req.Description,
		Goal:          req.Goal,
		Progress:      req.Progress,
		PointsAwarded: req.PointsAwarded,
		Status:        req.Status,
	}

	if req.Deadline != nil {
		t, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido, use YYYY-MM-DD"})
			return
		}
		reto.Deadline = &t
	}

	err = ctrl.useCase.Execute(&reto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reto actualizado exitosamente",
	})
}
