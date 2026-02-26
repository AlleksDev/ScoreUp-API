package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/AlleksDev/ScoreUp-API/internal/user/application"
	"github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure/controllers"
	"github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure/repository"
)

type UserModule struct {
	CreateCtrl *controllers.CreateUserController
}

// NewUserModule es el provider que Wire usa para construir el módulo.
func NewUserModule(
	create *controllers.CreateUserController,
) *UserModule {
	return &UserModule{
		CreateCtrl: create,
	}
}

// RegisterRoutes registra las rutas del feature User en el gin.Engine.
func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/api/users")
	{
		users.POST("", m.CreateCtrl.Handle)
	}
}

// UserProviderSet es el conjunto de providers que Wire inyecta automáticamente.
// Wire analiza los parámetros de entrada/salida de cada constructor y resuelve
// el grafo de dependencias en TIEMPO DE COMPILACIÓN (como Hilt + KSP en Kotlin).
var UserProviderSet = wire.NewSet(
	repository.NewUserMySQLRepository,
	app.NewCreateUser,
	controllers.NewCreateUserController,
	NewUserModule,
)