package controllers

import (
	"log"
	"net/http"
	"runtime/debug"
	"strconv"

	app "github.com/AlleksDev/ScoreUp-API/internal/reto/application"
	"github.com/AlleksDev/ScoreUp-API/internal/websocket"
	"github.com/gin-gonic/gin"
)

type DeleteRetoController struct {
	useCase    *app.DeleteReto
	getUseCase *app.GetReto
	hub        *websocket.Hub
}

func NewDeleteRetoController(useCase *app.DeleteReto, getReto *app.GetReto, hub *websocket.Hub) *DeleteRetoController {
	return &DeleteRetoController{useCase: useCase, getUseCase: getReto, hub: hub}
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

	// Push: notificar a clientes WS suscritos al canal "retos".
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC][broadcastRetos] recovered: %v\n%s", r, debug.Stack())
			}
		}()
		ctrl.broadcastRetos()
	}()
}

func (ctrl *DeleteRetoController) broadcastRetos() {
	retos, err := ctrl.getUseCase.ExecuteAll()
	if err != nil {
		log.Printf("[WS] Error obteniendo retos para broadcast: %v", err)
		return
	}
	ctrl.hub.BroadcastJSON("retos", gin.H{"retos": retos})
}
