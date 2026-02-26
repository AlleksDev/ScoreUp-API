package controllers

import (
	"net/http"
	"time"

	app "github.com/AlleksDev/ScoreUp-API/internal/reto/application"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/domain/entities"
	"github.com/gin-gonic/gin"
)

type CreateRetoController struct {
	useCase *app.CreateReto
}

func NewCreateRetoController(useCase *app.CreateReto) *CreateRetoController {
	return &CreateRetoController{useCase: useCase}
}

type createRetoRequest struct {
	Subject       string  `json:"subject" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	Goal          int     `json:"goal" binding:"required,min=1"`
	PointsAwarded int     `json:"points_awarded"`
	Deadline      *string `json:"deadline"`
}

func (ctrl *CreateRetoController) Handle(c *gin.Context) {
	var req createRetoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	reto := entities.Reto{
		UserID:        userID.(int64),
		Subject:       req.Subject,
		Description:   req.Description,
		Goal:          req.Goal,
		Progress:      0,
		PointsAwarded: req.PointsAwarded,
		Status:        "activo",
	}

	if reto.PointsAwarded == 0 {
		reto.PointsAwarded = 20
	}

	if req.Deadline != nil {
		t, err := time.Parse("2006-01-02", *req.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de fecha inválido, use YYYY-MM-DD"})
			return
		}
		reto.Deadline = &t
	}

	err := ctrl.useCase.Execute(&reto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reto creado exitosamente",
		"id":      reto.ID,
	})
}
