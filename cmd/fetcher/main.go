package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/dmitryshur/hackernews/internal/data"
	"github.com/dmitryshur/hackernews/internal/jsonlog"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"time"
)

const baseUrl = "https://hacker-news.firebaseio.com/v0"

// TODO: need a method to update a story (fetch all the comments). when the user re-enters a story, some newer
//  comments might be missing
func main() {
	var cfg config

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	cfg.fetchInterval = time.Second * 600

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	api := NewApi(http.DefaultClient, baseUrl)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	s := data.NewModel(db)
	f := NewFetcher(cfg, logger, api, s)

	f.Start()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
