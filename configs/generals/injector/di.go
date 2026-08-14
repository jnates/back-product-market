// Package injector builds the application's dependency injection container.
package injector

import (
	"backend_crudgo/configs"
	"backend_crudgo/configs/generals/router"
	"backend_crudgo/configs/generals/router/health"
	productsRoutes "backend_crudgo/configs/generals/router/products"
	usersRoutes "backend_crudgo/configs/generals/router/users"
	"backend_crudgo/configs/storage"
	"backend_crudgo/pkg/middleware"
	productsHandler "backend_crudgo/products/handler"
	productsRepository "backend_crudgo/products/repository"
	productsService "backend_crudgo/products/service"
	usersHandler "backend_crudgo/users/handler"
	usersRepository "backend_crudgo/users/repository"
	usersService "backend_crudgo/users/service"

	customserverEcho "github.com/jnates/go-toolkit/tools/customserver/echo"

	"go.uber.org/dig"
)

// BuildContainer wires every dependency the application needs via dig.
// Panics if a provider is misconfigured, since that is a programming error
// that should fail fast at startup rather than surface at request time.
func BuildContainer(config *configs.Config) *dig.Container {
	container := dig.New()

	provide(container, func() *configs.Config { return config })
	provide(container, storage.PostgresConnection)
	provide(container, customserverEcho.NewServer)
	provide(container, middleware.NewAuthMiddleware)

	provide(container, productsRepository.NewProductRepository)
	provide(container, productsService.NewProductService)
	provide(container, productsHandler.NewProductHandler)
	provide(container, productsRoutes.NewRouteProducts)

	provide(container, usersRepository.NewUserRepository)
	provide(container, usersService.NewUserService)
	provide(container, usersHandler.NewUserHandler)
	provide(container, usersRoutes.NewRouteUsers)

	provide(container, health.NewHandler)
	provide(container, router.NewRouter)

	return container
}

// provide registers a constructor and panics if dig rejects it (duplicate or malformed provider).
func provide(container *dig.Container, constructor interface{}) {
	if err := container.Provide(constructor); err != nil {
		panic(err)
	}
}
