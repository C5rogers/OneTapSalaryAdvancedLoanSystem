package main

import (
	"net/http"
	"strings"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/api/server"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/cloudinary"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/db"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/lib"
	"github.com/rs/cors"
)

func run(configPath *string) error {
	// load configurations
	config, err := config.Load(configPath)

	if err != nil {
		return err
	}

	// setup go-database connection
	dbConn, err := db.NewClient(config)
	if err != nil {
		return err
	}

	// setup cloudinary client
	cloudinaryClient, err := cloudinary.NewCloudinaryClient(config)
	if err != nil {
		return nil
	}

	// setup logger
	logger := lib.NewLogger(config.Log)

	// boot up progress text
	logger.With("config", config).Debug("configuration loaded")

	// no need to provide poller we just want to use the APIs only
	s, err := server.NewServer(logger, config, cloudinaryClient, dbConn)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	s.ApplyRoutes(mux)

	// apply middleware to extract the  session user information from authorization header
	wrappedMux := s.ExtractUserFromToken(mux)

	// boot up progress text
	logger.Info("server routes applied")

	allowedOrigins := strings.Split(config.Server.AllowedOrigins, ",")

	if len(allowedOrigins) == 0 {
		allowedOrigins = append(allowedOrigins, "*")
	}

	// TODO: only allow applicable methods and headers
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"OPTIONS", "GET", "HEAD", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	s.Handler = c.Handler(wrappedMux)

	// boot up progress text
	logger.Info("server is listening", "PORT", s.Addr)

	return s.ListenAndServe()
}
