package controllers

import (
	"net/http"

	app "github.com/AlleksDev/ScoreUp-API/internal/user/application"
	"github.com/gin-gonic/gin"
)

type GetRankController struct {
	useCase *app.GetRank
}

func NewGetRankController(useCase *app.GetRank) *GetRankController {
	return &GetRankController{useCase: useCase}
}

func (ctrl *GetRankController) Handle(c *gin.Context) {

	rank, err := ctrl.useCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error obteniendo ranking",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ranking": rank})
}
