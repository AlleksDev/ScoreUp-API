//go:build wireinject
// +build wireinject

// wire.go — Este archivo es el INJECTOR de Google Wire.
// Define QUÉ providers usar. Wire genera el código de cableado en wire_gen.go.

package main

import (
	"github.com/google/wire"

	"github.com/AlleksDev/ScoreUp-API/internal/core"
	logroInfra "github.com/AlleksDev/ScoreUp-API/internal/logro/infrastructure"
	retoInfra "github.com/AlleksDev/ScoreUp-API/internal/reto/infrastructure"
	userInfra "github.com/AlleksDev/ScoreUp-API/internal/user/infrastructure"
	usuarioLogroInfra "github.com/AlleksDev/ScoreUp-API/internal/usuario_logro/infrastructure"
	usuarioRetoInfra "github.com/AlleksDev/ScoreUp-API/internal/usuario_reto/infrastructure"
)

// InitializeApp le dice a Wire: "construí un *App usando TODOS estos providers".
// Wire analiza el grafo de dependencias y genera wire_gen.go con el código real.
func InitializeApp() (*App, error) {
	wire.Build(
		// Infraestructura base
		core.GetMySQLPool, // → *core.Conn_MySQL

		// Feature: User (repo → adapters → use cases → controllers → module)
		userInfra.UserProviderSet,

		// Feature: Reto
		retoInfra.RetoProviderSet,

		// Feature: Logro
		logroInfra.LogroProviderSet,

		// Feature: UsuarioReto (relación M:N usuario-reto)
		usuarioRetoInfra.UsuarioRetoProviderSet,

		// Feature: UsuarioLogro (evaluación de logros)
		usuarioLogroInfra.UsuarioLogroProviderSet,

		// Composición final
		ProvideGinEngine, // → *gin.Engine
		NewApp,           // → *App
	)
	return nil, nil
}
