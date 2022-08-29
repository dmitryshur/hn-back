package server

import (
	"github.com/dmitryshur/hackernews/internal/jsonlog"
	"sync"
)

type Config struct {
	Port    int
	Env     string
	Limiter struct {
		Rps     float64
		Burst   int
		Enabled bool
	}
	Cors struct {
		TrustedOrigins []string
	}
}

type application struct {
	config Config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
}

func NewApplication(cfg Config, logger *jsonlog.Logger) *application {
	return &application{config: cfg, logger: logger}
}
