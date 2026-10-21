package registry

import (
	"github.com/ryanpujo/pendataan-apar-backend/internal/controllers"
	"github.com/ryanpujo/pendataan-apar-backend/internal/repositories"
	"github.com/ryanpujo/pendataan-apar-backend/internal/services"
)

func (r *registry) NewCredentialRepo() *repositories.CredentialRepo {
	return repositories.NewCredentialRepo(r.db)
}

func (r *registry) NewCredentialService() *services.CredentialService {
	return services.NewCredentialService(r.NewCredentialRepo())
}

func (r *registry) NewCredentialController() *controllers.CredentialController {
	return controllers.NewCredentialController(r.NewCredentialService())
}
