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
			log.Println("Error create config: " + err.Error())
			return
		}
		appLog, err := logger.New(conf.Logger)
		if err != nil {
			log.Println("Error create app logger: " + err.Error())
			return
		}
		serverLog, err := logger.New(conf.GrpcLogger)
		if err != nil {
			log.Println("Error create server logger: " + err.Error())
			return
		}

		s, err := storage.CreateStorage(conf.Database)
		if err != nil {
			log.Println("Error create storage: " + err.Error())
			return
		}

		server := grpc.NewServer(serverLog, app.New(appLog, s), conf.GrpcServer)
		go func() {
			if err := server.Listen(ctx); err != nil {
				log.Println("Error execute app: " + err.Error())
				cancel()
			}
		}()
		<-ctx.Done()
	},
}
