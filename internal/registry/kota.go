package registry

import (
	"github.com/ryanpujo/pendataan-apar-backend/internal/controllers"
	"github.com/ryanpujo/pendataan-apar-backend/internal/repositories"
	"github.com/ryanpujo/pendataan-apar-backend/internal/services"
)

func (r *registry) NewKotaRepo() *repositories.KotaRepo {
	return repositories.NewKotaRepo(r.db)
}

func (r *registry) NewKotaService() *services.KotaService {
	return services.NewKotaService(r.NewKotaRepo())
}

func (r *registry) NewKotaController() *controllers.KotaController {
	return controllers.NewKotaController(r.NewKotaService())
}