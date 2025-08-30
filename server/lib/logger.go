package lib

import (
	"log/slog"
	"os"

	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
)

func NewLogger(cfg config.Log) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Level,
	}))
}
