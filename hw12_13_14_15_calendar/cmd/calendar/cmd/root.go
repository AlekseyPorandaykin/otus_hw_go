package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/app"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/server/http"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/spf13/cobra"
)

var configFile string

var RootCmd = &cobra.Command{
	Short: "Calendar app",
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

		serverLog, err := logger.New(conf.HTTPLogger)
		if err != nil {
			log.Println("Error create server logger: " + err.Error())
			return
		}

		s, err := storage.CreateStorage(conf.Database)
		if err != nil {
			log.Println("Error create storage: " + err.Error())
			return
		}

		server := internalhttp.NewServer(serverLog, app.New(appLog, s), conf.HTTPServer)
		go func() {
			if err := server.Start(); err != nil {
				log.Println("failed to start http server: " + err.Error())
				cancel()
			}
		}()
		defer server.Stop()

		<-ctx.Done()
	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(grpcServerCmd)
	RootCmd.AddCommand(schedulerCmd)
	RootCmd.AddCommand(senderCmd)
}
