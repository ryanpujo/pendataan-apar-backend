package routes

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/ryanpujo/pendataan-apar-backend/internal/controllers"
)

type ControllersRegistry struct {
	Credential controllers.Credential
	Kota       controllers.Kota
}

func Route(registry *ControllersRegistry) *gin.Engine {
	r := gin.Default()
	credentialRoute(r, registry.Credential)
	kotaRoute(r, registry.Kota)

	return r
}

func credentialRoute(r *gin.Engine, credController controllers.Credential) {
	v1 := "/v1"
	v1CredentialsRoute := r.Group(fmt.Sprintf("/%s/credentials", v1))
	v1CredentialsRoute.POST("/register", credController.Register)
	v1CredentialsRoute.GET("/:id", credController.FindByID)
}

func kotaRoute(r *gin.Engine, kotaController controllers.Kota) {
	v1Route := r.Group("/v1/kota")
	v1Route.POST("/create", kotaController.Create)
	v1Route.GET("/:id", kotaController.FindByID)
}
