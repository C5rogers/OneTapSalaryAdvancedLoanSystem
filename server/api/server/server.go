package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/cloudinary"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/db"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/security"
)

type Server struct {
	http.Server
	Logger      *slog.Logger
	Cloudinary  *cloudinary.CloudinaryClient
	Config      *config.Config
	DB          *db.Database
	RateLimiter *security.RateLimiter
}

func NewServer(logger *slog.Logger, cfg *config.Config, Cloudinary *cloudinary.CloudinaryClient, database *db.Database) (*Server, error) {
	return &Server{
		Server: http.Server{
			Addr: cfg.Server.ListenAddress,
		},
		Logger:      logger,
		Cloudinary:  Cloudinary,
		Config:      cfg,
		DB:          database,
		RateLimiter: security.NewRateLimiter(5, time.Minute), // 5 attempts / 1 min
	}, nil
}
