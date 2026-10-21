package controllers_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/internal/controllers"
	"github.com/ryanpujo/pendataan-apar-backend/internal/routes"
)

var (
	sm         *ServiceMock
	ksm *KotaServiceMock
	controller controllers.Credential
	kotaController controllers.Kota
	handler    http.Handler
	v1         = "/v1"
	now        = time.Now().Truncate(0)
)

func TestMain(m *testing.M) {
	sm = new(ServiceMock)
	ksm = new(KotaServiceMock)
	controller = controllers.NewCredentialController(sm)
	kotaController = controllers.NewKotaController(ksm)
	controllerRegistry := routes.ControllersRegistry{
		Credential: controller,
		Kota: kotaController,
	}
	handler = routes.Route(&controllerRegistry)
	os.Exit(m.Run())
}
