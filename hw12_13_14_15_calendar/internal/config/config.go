package config

import (
	internalhttp "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/config"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
)

type Config struct {
	Logger     *logger.Config       `mapstructure:"logger"`
	HTTPLogger *logger.Config       `mapstructure:"http_logger"`
	Database   *StorageConfig       `mapstructure:"database"`
	HTTPServer *internalhttp.Config `mapstructure:"http"`
}
type StorageConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Storage  string `mapstructure:"storage"`
	Driver   string `mapstructure:"driver"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"db_name"`
	SslMode  string `mapstructure:"ssl_mode"`
	Password string `mapstructure:"password"`
}

func New(pathToFile string) (Config, error) {
	conf, err := config.CreateConfig(pathToFile, "toml", Config{})
	return conf.(Config), err
}
