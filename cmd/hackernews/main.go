package main

import (
	"flag"
	"os"
	"strings"
)
import "github.com/dmitryshur/hackernews/pkg/jsonlog"
import "github.com/dmitryshur/hackernews/pkg/server"

// TODO: add type param (newest, best)
// TODO: add pagination params (offset, limit)
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

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	app := server.NewApplication(cfg, logger)
	err := app.Serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
