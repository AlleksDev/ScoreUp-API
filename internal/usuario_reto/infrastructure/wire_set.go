package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	app "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/application"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/ports"
	domainRepo "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/domain/repository"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/infrastructure/adapters"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/infrastructure/controllers"
	"github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/infrastructure/repository"
)

type UsuarioRetoModule struct {
	JoinCtrl     *controllers.JoinRetoController
	ProgressCtrl *controllers.UpdateProgressController
	GetCtrl      *controllers.GetUsuarioRetosController
	LeaveCtrl    *controllers.LeaveRetoController
}

func NewUsuarioRetoModule(
	join *controllers.JoinRetoController,
	progress *controllers.UpdateProgressController,
	get *controllers.GetUsuarioRetosController,
	leave *controllers.LeaveRetoController,
) *UsuarioRetoModule {
	return &UsuarioRetoModule{
		JoinCtrl:     join,
		ProgressCtrl: progress,
		GetCtrl:      get,
		LeaveCtrl:    leave,
	}
}

func (m *UsuarioRetoModule) RegisterRoutes(r *gin.RouterGroup) {
	ur := r.Group("/usuario-retos")
	{
		ur.POST("", m.JoinCtrl.Handle)
		ur.POST("/:retoId/join", m.JoinCtrl.HandleByParam)
		ur.PUT("/:retoId/progress", m.ProgressCtrl.Handle)
		ur.GET("", m.GetCtrl.HandleByUser)
		ur.GET("/:retoId/participants", m.GetCtrl.HandleByReto)
		ur.DELETE("/:retoId", m.LeaveCtrl.Handle)
	}
}

var UsuarioRetoProviderSet = wire.NewSet(
	// Repository
	repository.NewUsuarioRetoMySQLRepository,
	wire.Bind(new(domainRepo.UsuarioRetoRepository), new(*repository.UsuarioRetoMySQLRepository)),

	// Adapters (implementan los puertos para interactuar con otros módulos)
	adapters.NewUserScoreAdapter,
	wire.Bind(new(ports.UserScorePort), new(*adapters.UserScoreAdapter)),

	adapters.NewRetoQueryAdapter,
	wire.Bind(new(ports.RetoQueryPort), new(*adapters.RetoQueryAdapter)),

	adapters.NewLogroEvaluatorAdapter,
	wire.Bind(new(ports.LogroEvaluatorPort), new(*adapters.LogroEvaluatorAdapter)),

	// Use Cases
	app.NewJoinReto,
	app.NewUpdateProgress,
	app.NewGetUsuarioRetos,
	app.NewLeaveReto,

	// Controllers
	controllers.NewJoinRetoController,
	controllers.NewUpdateProgressController,
	controllers.NewGetUsuarioRetosController,
	controllers.NewLeaveRetoController,

	// Module
	NewUsuarioRetoModule,
)
