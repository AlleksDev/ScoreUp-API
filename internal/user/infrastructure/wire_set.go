package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	app "github.com/AlleksDev/ScoreUp-API/internal/user/application"
	"github.com/AlleksDev/ScoreUp-API/internal/user/domain/ports"
	domainRepo "github.com/AlleksDev/ScoreUp-API/internal/user/domain/repository"
	"github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure/adapters"
	"github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure/controllers"
	"github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure/repository"
)

type UserModule struct {
	CreateCtrl *controllers.CreateUserController
	LoginCtrl  *controllers.LoginUserController
	RankCtrl   *controllers.GetRankController
}

func NewUserModule(
	create *controllers.CreateUserController,
	login *controllers.LoginUserController,
	rank *controllers.GetRankController,
) *UserModule {
	return &UserModule{
		CreateCtrl: create,
		LoginCtrl:  login,
		RankCtrl:   rank,
	}
}

// RegisterRoutes registra las rutas públicas del feature User (sin auth).
func (m *UserModule) RegisterRoutes(r *gin.Engine) {
	users := r.Group("/api/users")
	{
		users.POST("/register", m.CreateCtrl.Handle)
		users.POST("/login", m.LoginCtrl.Handle)
		users.GET("/rank", m.RankCtrl.Handle)
	}
}

var UserProviderSet = wire.NewSet(
	// Repository
	repository.NewUserMySQLRepository,
	wire.Bind(new(domainRepo.UserRepository), new(*repository.UserMySQLRepository)),

	// Adapters
	adapters.NewBcryptAdapter,
	wire.Bind(new(ports.IBcryptService), new(*adapters.BcryptAdapter)),

	adapters.NewJWTTokenAdapter,
	wire.Bind(new(ports.TokenManager), new(*adapters.JWTTokenAdapter)),

	// Use Cases
	app.NewCreateUser,
	app.NewLoginUser,
	app.NewGetRank,

	// Controllers
	controllers.NewCreateUserController,
	controllers.NewLoginUserController,
	controllers.NewGetRankController,

	// Module
	NewUserModule,
)
