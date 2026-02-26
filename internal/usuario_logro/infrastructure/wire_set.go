package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/application"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/ports"
	domainRepo "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/domain/repository"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/infrastructure/adapters"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/infrastructure/controllers"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/infrastructure/repository"
)

type UsuarioLogroModule struct {
	EvaluateCtrl *controllers.EvaluateLogrosController
	GetCtrl      *controllers.GetUsuarioLogrosController
	DeleteCtrl   *controllers.DeleteUsuarioLogroController
}

func NewUsuarioLogroModule(
	evaluate *controllers.EvaluateLogrosController,
	get *controllers.GetUsuarioLogrosController,
	del *controllers.DeleteUsuarioLogroController,
) *UsuarioLogroModule {
	return &UsuarioLogroModule{
		EvaluateCtrl: evaluate,
		GetCtrl:      get,
		DeleteCtrl:   del,
	}
}

func (m *UsuarioLogroModule) RegisterRoutes(r *gin.RouterGroup) {
	ul := r.Group("/usuario-logros")
	{
		ul.POST("/evaluate", m.EvaluateCtrl.Handle)
		ul.GET("", m.GetCtrl.Handle)
		ul.DELETE("/:logroId", m.DeleteCtrl.Handle)
	}
}

var UsuarioLogroProviderSet = wire.NewSet(
	// Repository
	repository.NewUsuarioLogroMySQLRepository,
	wire.Bind(new(domainRepo.UsuarioLogroRepository), new(*repository.UsuarioLogroMySQLRepository)),

	// Adapters (implementan los puertos para consultar otros módulos)
	adapters.NewUserQueryAdapter,
	wire.Bind(new(ports.UserQueryPort), new(*adapters.UserQueryAdapter)),

	adapters.NewRetoQueryAdapter,
	wire.Bind(new(ports.RetoQueryPort), new(*adapters.RetoQueryAdapter)),

	adapters.NewLogroQueryAdapter,
	wire.Bind(new(ports.LogroQueryPort), new(*adapters.LogroQueryAdapter)),

	// Use Cases
	app.NewEvaluateLogros,
	app.NewGetUsuarioLogros,
	app.NewDeleteUsuarioLogro,

	// Controllers
	controllers.NewEvaluateLogrosController,
	controllers.NewGetUsuarioLogrosController,
	controllers.NewDeleteUsuarioLogroController,

	// Module
	NewUsuarioLogroModule,
)
