package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/app"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/spf13/cobra"
)

var grpcServerCmd = &cobra.Command{
	Use:   "grpc_server",
	Short: "Start grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()
		conf, err := config.New(configFile)
		if err != nil {
			log.Panic("Error create config: " + err.Error())
		}
		appLog, err := logger.New(conf.Logger)
		if err != nil {
			log.Panic("Error create app logger: " + err.Error())
		}
		serverLog, err := logger.New(conf.GrpcLogger)
		if err != nil {
			log.Panic("Error create server logger: " + err.Error())
		}

		s, err := storage.CreateStorage(conf.Database)
		if err != nil {
			log.Panic("Error create storage: " + err.Error())
		}

		server := grpc.NewServer(serverLog, app.New(appLog, s), conf.GrpcServer)
		if err := server.Listen(ctx); err != nil {
			log.Fatal(err)
		}
	},
}
