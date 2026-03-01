package controllers

import (
	"net/http"
	"strconv"

	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/application"
	"github.com/gin-gonic/gin"
)

type LeaveRetoController struct {
	useCase *app.LeaveReto
}

func NewLeaveRetoController(useCase *app.LeaveReto) *LeaveRetoController {
	return &LeaveRetoController{useCase: useCase}
}

func (ctrl *LeaveRetoController) Handle(c *gin.Context) {
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

	err = ctrl.useCase.Execute(uid, retoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Abandonaste el reto exitosamente",
	})
}
