package main

import (
	"flag"
	"github.com/dmitryshur/hackernews/internal/jsonlog"
	"github.com/dmitryshur/hackernews/internal/server"
	"os"
	"strings"
)

func main() {
	var cfg server.Config

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")

	flag.Float64Var(&cfg.Limiter.Rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.Limiter.Burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.Limiter.Enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.Cors.TrustedOrigins = strings.Fields(val)
		return nil
	})

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	app := server.NewApplication(cfg, logger)
	err := app.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
