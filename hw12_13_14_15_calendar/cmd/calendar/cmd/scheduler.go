package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue/ampq"
	"github.com/spf13/cobra"
)

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Start scheduler for send remind notification",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()
		conf, err := config.NewSchedulerApp(configFile)
		if err != nil {
			log.Panic("Error create config: ", err)
		}
		appLog, err := logger.New(conf.Logger)
		if err != nil {
			log.Panic("Error create app logger: ", err)
		}
		sender := ampq.NewProducer(ampq.NewConnection(conf.Producer, appLog), appLog)
		db, err := storage.CreateStorage(conf.Database)
		if err != nil {
			log.Panic("Error create storage: " + err.Error())
		}
		app := scheduler.New(appLog, db, sender, conf.Scheduler)

		go func() {
			if err := app.Run(ctx); err != nil {
				log.Println(err)
				cancel()
			}
		}()

		<-ctx.Done()
	},
}
