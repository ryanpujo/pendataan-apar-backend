package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ryanpujo/pendataan-apar-backend/config"
	"github.com/ryanpujo/pendataan-apar-backend/internal/database"
	"github.com/ryanpujo/pendataan-apar-backend/internal/registry"
	"github.com/ryanpujo/pendataan-apar-backend/internal/routes"
)

func main() {
	DB, err := database.NewPostgresDB()
	if err != nil {
		panic(fmt.Sprintf("failed to establish database connection: %v", err))
	}

	r := registry.NewRegistry(DB)
	controllers := routes.ControllersRegistry{
		Credential: r.NewCredentialController(),
		Kota: r.NewKotaController(),
	}
	if err := application(routes.Route(&controllers)); err != nil {
		panic("failed to run the entire backend")
	}
}

func application(handler http.Handler) error {
	server := http.Server{
		Addr:              fmt.Sprintf(":%s", config.GetConfig().AppPort),
		Handler:           handler,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return server.ListenAndServe()
}
