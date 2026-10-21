package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ryanpujo/pendataan-apar-backend/config"
)

func NewPostgresDB() (*sql.DB, error) {
	cfg := config.GetConfig()
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Konfigurasi Connection Pool
	db.SetMaxOpenConns(25)                 // Maksimal koneksi aktif bersamaan
	db.SetMaxIdleConns(10)                 // Koneksi siaga di pool
	db.SetConnMaxLifetime(5 * time.Minute) // Umur maksimal koneksi sebelum di-recycle

	// Pastikan database benar-benar merespons (Ping)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}