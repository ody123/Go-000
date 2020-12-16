// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"demo/internal/dao"
	"demo/internal/server/grpc"
	"demo/internal/server/http"
	"demo/internal/service"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp))
}
