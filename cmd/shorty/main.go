package main

import (
	"context"
	"os"
	"strings"
	"time"

	"git.6summits.net/srv/shorty/internals/shorty"
	"git.6summits.net/srv/shorty/pkg/mongodb"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HTTPAddress string `env:"HTTP_ADDRESS" envDefault:":8080"`
	MongoDBURI  string `env:"MONGO_CONNECTION_STRING,notEmpty"`
	MongoDBName string `env:"MONGO_NAME" envDefault:"shorty"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat   string `env:"LOG_FORMAT" envDefault:"json"`
}

const (
	gracePeriod = 30 * time.Second
)

func main() {
	if err := r(); err != nil {
		log.Error().Err(err).Msg("")
		os.Exit(1)
	}
}

func r() error {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return err
	}
	log.Info().Interface("Config", cfg).Msg("")

	l, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		return err
	}

	zerolog.SetGlobalLevel(l)
	if strings.ToLower(cfg.LogFormat) == "console" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
	defer cancel()

	mongoClient, err := mongodb.New(ctx, cfg.MongoDBURI, cfg.MongoDBName)
	if err != nil {
		return err
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Error().Err(err).Msg("")
		}
	}()

	s := shorty.New(mongoClient)

	return s.ListenAndServe(cfg.HTTPAddress)
}
