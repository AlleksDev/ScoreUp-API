package controllers

import (
	"net/http"
	"strconv"

	"github.com/AlleksDev/ScoreUp-API/internal/middleware"
	app "github.com/AlleksDev/ScoreUp-API/internal/reto/application"
	"github.com/gin-gonic/gin"
)

type GetRetoController struct {
	useCase *app.GetReto
}

func NewGetRetoController(useCase *app.GetReto) *GetRetoController {
	return &GetRetoController{useCase: useCase}
}

func (ctrl *GetRetoController) HandleGetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	reto, err := ctrl.useCase.ExecuteByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"reto": reto,
	})
}

func (ctrl *GetRetoController) HandleGetAll(c *gin.Context) {
	retos, err := ctrl.useCase.ExecuteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	payload := gin.H{"retos": retos}

	c.JSON(http.StatusOK, payload)
}

func (ctrl *GetRetoController) HandleGetByCreator(c *gin.Context) {
	uid, ok := middleware.GetUserID(c)
	if !ok {
		return
	}

	retos, err := ctrl.useCase.ExecuteByCreator(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"retos": retos,
	})
}
