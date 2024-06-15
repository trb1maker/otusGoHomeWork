package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/go-yaml/yaml"
)

type config struct {
	Rabbit   rabbit        `yaml:"rabbit"`
	Postgres postgres      `yaml:"postgres"`
	Logger   logger        `yaml:"logger"`
	Interval time.Duration `yaml:"interval"`
}

type rabbit struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type postgres struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type logger struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func loadConfig(fileName string) (*config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var conf config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	conf.Logger.init()

	return &conf, nil
}

func (l *logger) init() {
	var options slog.HandlerOptions

	switch l.Level {
	case "debug":
		options.AddSource = true
		options.Level = slog.LevelDebug
	case "info":
		options.Level = slog.LevelInfo
	case "warn":
		options.Level = slog.LevelWarn
	case "error":
		options.AddSource = true
		options.Level = slog.LevelError
	default:
		options.Level = slog.LevelError
	}

	if l.Format == "json" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(
			os.Stdout,
			&options,
		)))
		return
	}

	slog.SetDefault(slog.New(slog.NewTextHandler(
		os.Stdout,
		&options,
	)))
}
