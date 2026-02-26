package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	app "github.com/AlleksDev/ScoreUp-API/internal/logro/application"
	domainRepo "github.com/AlleksDev/ScoreUp-API/internal/logro/domain/repository"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/infrastructure/controllers"
	"github.com/AlleksDev/ScoreUp-API/internal/logro/infrastructure/repository"
)

type LogroModule struct {
	CreateCtrl *controllers.CreateLogroController
	GetCtrl    *controllers.GetLogroController
	UpdateCtrl *controllers.UpdateLogroController
	DeleteCtrl *controllers.DeleteLogroController
}

func NewLogroModule(
	create *controllers.CreateLogroController,
	get *controllers.GetLogroController,
	update *controllers.UpdateLogroController,
	del *controllers.DeleteLogroController,
) *LogroModule {
	return &LogroModule{
		CreateCtrl: create,
		GetCtrl:    get,
		UpdateCtrl: update,
		DeleteCtrl: del,
	}
}

func (m *LogroModule) RegisterRoutes(r *gin.RouterGroup) {
	logros := r.Group("/logros")
	{
		logros.POST("", m.CreateCtrl.Handle)
		logros.GET("", m.GetCtrl.HandleGetAll)
		logros.GET("/:id", m.GetCtrl.HandleGetByID)
		logros.PUT("/:id", m.UpdateCtrl.Handle)
		logros.DELETE("/:id", m.DeleteCtrl.Handle)
	}
}

var LogroProviderSet = wire.NewSet(
	repository.NewLogroMySQLRepository,
	wire.Bind(new(domainRepo.LogroRepository), new(*repository.LogroMySQLRepository)),
	app.NewCreateLogro,
	app.NewGetLogro,
	app.NewUpdateLogro,
	app.NewDeleteLogro,
	controllers.NewCreateLogroController,
	controllers.NewGetLogroController,
	controllers.NewUpdateLogroController,
	controllers.NewDeleteLogroController,
	NewLogroModule,
)
