package config

import (
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/config"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue"
)

type Config struct {
	Logger     *logger.Config       `mapstructure:"logger"`
	HTTPLogger *logger.Config       `mapstructure:"http_logger"`
	GrpcLogger *logger.Config       `mapstructure:"grpc_logger"`
	Database   *StorageConfig       `mapstructure:"database"`
	HTTPServer *internalhttp.Config `mapstructure:"http"`
	GrpcServer *grpc.Config         `mapstructure:"grpc"`
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

type SchedulerApp struct {
	Logger    *logger.Config    `mapstructure:"logger"`
	Database  *StorageConfig    `mapstructure:"database"`
	Scheduler *scheduler.Config `mapstructure:"scheduler"`
	Producer  *queue.Config     `mapstructure:"producer"`
}

type SenderApp struct {
	Logger   *logger.Config `mapstructure:"logger"`
	Database *StorageConfig `mapstructure:"database"`
	Consumer *queue.Config  `mapstructure:"consumer"`
}

func New(pathToFile string) (Config, error) {
	conf, err := config.CreateConfig(pathToFile, "toml", Config{})
	return conf.(Config), err
}

func NewSchedulerApp(pathToFile string) (SchedulerApp, error) {
	conf, err := config.CreateConfig(pathToFile, "toml", SchedulerApp{})
	return conf.(SchedulerApp), err
}

func NewSenderApp(pathToFile string) (SenderApp, error) {
	conf, err := config.CreateConfig(pathToFile, "toml", SenderApp{})
	return conf.(SenderApp), err
}
