package controllers

import (
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
	userApp "github.com/AlleksDev/ScoreUp-API/internal/user/application"
	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/application"
	"github.com/AlleksDev/ScoreUp-API/internal/websocket"
	"github.com/gin-gonic/gin"
)

type UpdateProgressController struct {
	useCase   *app.UpdateProgress
	getRankUC *userApp.GetRank
	hub       *websocket.Hub
}

func NewUpdateProgressController(useCase *app.UpdateProgress, getRank *userApp.GetRank, hub *websocket.Hub) *UpdateProgressController {
	return &UpdateProgressController{useCase: useCase, getRankUC: getRank, hub: hub}
}

type updateProgressRequest struct {
	Progress int `json:"progress" binding:"required,min=0"`
}

func (ctrl *UpdateProgressController) Handle(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	retoIDParam := c.Param("retoId")
	retoID, err := strconv.ParseInt(retoIDParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reto inválido"})
		return
	}

	var req updateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	result, err := ctrl.useCase.Execute(uid, retoID, req.Progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"message":   "Progreso actualizado exitosamente",
		"completed": result.Completed,
	}

	if len(result.LogrosAwarded) > 0 {
		response["logros_awarded"] = result.LogrosAwarded
	}

	c.JSON(http.StatusOK, response)

	// Push: notificar ranking actualizado a clientes WS del canal "rank".
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC][broadcastRank] recovered: %v\n%s", r, debug.Stack())
			}
		}()
		ctrl.broadcastRank()
	}()
}

func (ctrl *UpdateProgressController) broadcastRank() {
	rank, err := ctrl.getRankUC.Execute()
	if err != nil {
		log.Printf("[WS] Error obteniendo ranking para broadcast: %v", err)
		return
	}
	ctrl.hub.BroadcastJSON("rank", gin.H{"ranking": rank})
}
