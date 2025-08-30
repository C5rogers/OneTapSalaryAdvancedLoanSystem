package server

import (
	"log/slog"
	"net/http"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/cloudinary"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/db"
)

type Server struct {
	http.Server
	Logger     *slog.Logger
	Cloudinary *cloudinary.CloudinaryClient
	Config     *config.Config
	DB         *db.Database
}

func NewServer(logger *slog.Logger, cfg *config.Config, Cloudinary *cloudinary.CloudinaryClient, database *db.Database) (*Server, error) {
	return &Server{
		Server: http.Server{
			Addr: cfg.Server.ListenAddress,
		},
		Logger:     logger,
		Cloudinary: Cloudinary,
		Config:     cfg,
		DB:         database,
	}, nil
}
