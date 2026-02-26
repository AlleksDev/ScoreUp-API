package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	app "github.com/AlleksDev/ScoreUp-API/internal/reto/application"
	domainRepo "github.com/AlleksDev/ScoreUp-API/internal/reto/domain/repository"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/infrastructure/controllers"
	"github.com/AlleksDev/ScoreUp-API/internal/reto/infrastructure/repository"
)

type RetoModule struct {
	CreateCtrl *controllers.CreateRetoController
	GetCtrl    *controllers.GetRetoController
	UpdateCtrl *controllers.UpdateRetoController
	DeleteCtrl *controllers.DeleteRetoController
}

func NewRetoModule(
	create *controllers.CreateRetoController,
	get *controllers.GetRetoController,
	update *controllers.UpdateRetoController,
	del *controllers.DeleteRetoController,
) *RetoModule {
	return &RetoModule{
		CreateCtrl: create,
		GetCtrl:    get,
		UpdateCtrl: update,
		DeleteCtrl: del,
	}
}

func (m *RetoModule) RegisterRoutes(r *gin.RouterGroup) {
	retos := r.Group("/retos")
	{
		retos.POST("", m.CreateCtrl.Handle)
		retos.GET("", m.GetCtrl.HandleGetByUserID)
		retos.GET("/:id", m.GetCtrl.HandleGetByID)
		retos.PUT("/:id", m.UpdateCtrl.Handle)
		retos.DELETE("/:id", m.DeleteCtrl.Handle)
	}
}

var RetoProviderSet = wire.NewSet(
	repository.NewRetoMySQLRepository,
	wire.Bind(new(domainRepo.RetoRepository), new(*repository.RetoMySQLRepository)),
	app.NewCreateReto,
	app.NewGetReto,
	app.NewUpdateReto,
	app.NewDeleteReto,
	controllers.NewCreateRetoController,
	controllers.NewGetRetoController,
	controllers.NewUpdateRetoController,
	controllers.NewDeleteRetoController,
	NewRetoModule,
)
