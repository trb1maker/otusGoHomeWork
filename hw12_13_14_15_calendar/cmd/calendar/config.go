package main

import (
	"log/slog"
	"os"

	"github.com/go-yaml/yaml"
)

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

type config struct {
	Server  serverConf  `yaml:"server"`
	Storage storageConf `yaml:"storage"`
	Logger  loggerConf  `yaml:"logger"`
}

type serverConf struct {
	HTTP *httpConf `yaml:"http"`
}

type httpConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type storageConf struct {
	Type     string        `yaml:"type"`
	Postgres *postgresConf `yaml:"postgres"`
}

type postgresConf struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type loggerConf struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

func (l *loggerConf) init() {
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
