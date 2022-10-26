package server

import (
	"github.com/dmitryshur/hackernews/internal/data"
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
	Db struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
}

type application struct {
	config Config
	logger *jsonlog.Logger
	wg     sync.WaitGroup
	store  data.Db
}

func NewApplication(cfg Config, logger *jsonlog.Logger, store data.Db) *application {
	return &application{config: cfg, logger: logger, store: store}
}
