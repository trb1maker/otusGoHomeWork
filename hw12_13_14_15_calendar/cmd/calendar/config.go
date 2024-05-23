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

	if err := conf.init(); err != nil {
		return nil, err
	}

	return &conf, err
}

type config struct {
	Server  serverConf  `yaml:"server"`
	Storage storageConf `yaml:"storage"`
	Logger  loggerConf  `yaml:"logger"`
}

func (c *config) init() error {
	if err := c.Server.init(); err != nil {
		return err
	}

	if err := c.Storage.init(); err != nil {
		return err
	}

	c.Logger.init()
	return nil
}

type serverConf struct {
	HTTP *httpConf `yaml:"http"`
}

func (s *serverConf) init() error {
	if s.HTTP != nil {
		if err := s.HTTP.init(); err != nil {
			return err
		}
	}

	return nil
}

type httpConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (h *httpConf) init() error {
	if err := os.Setenv("SERVERHOST", h.Host); err != nil {
		return err
	}

	if err := os.Setenv("SERVERPORT", strconv.Itoa(h.Port)); err != nil {
		return err
	}

	return nil
}

type storageConf struct {
	Type     string        `yaml:"type"`
	Postgres *postgresConf `yaml:"postgres"`
}

func (s *storageConf) init() error {
	if s.Type == "postgres" {
		return s.Postgres.init()
	}
	return nil
}

type postgresConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (p *postgresConf) init() error {
	if err := os.Setenv("DBHOST", p.Host); err != nil {
		return err
	}

	if err := os.Setenv("DBPORT", strconv.Itoa(p.Port)); err != nil {
		return err
	}

	if err := os.Setenv("DBNAME", p.Database); err != nil {
		return err
	}

	if err := os.Setenv("DBUSER", p.User); err != nil {
		return err
	}

	if err := os.Setenv("DBPASSWORD", p.Password); err != nil {
		return err
	}

	return nil
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
