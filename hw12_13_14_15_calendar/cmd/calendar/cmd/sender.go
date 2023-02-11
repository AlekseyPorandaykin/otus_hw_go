package cmd

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/config"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/internal/sender"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/logger"
	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/pkg/queue/ampq"
	"github.com/spf13/cobra"
)

var senderCmd = &cobra.Command{
	Use:   "sender",
	Short: "Start sender for send remind notification",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()
		conf, err := config.NewSenderApp(configFile)
		if err != nil {
			log.Println("Error create config: ", err)
			return
		}
		appLog, err := logger.New(conf.Logger)
		if err != nil {
			log.Println("Error create app logger: ", err)
			return
		}
		con := ampq.NewConnection(conf.Consumer, appLog)
		if err != nil {
			log.Println("Error create storage: " + err.Error())
			return
		}
		consumer := ampq.NewConsumer(con, appLog)

		app := sender.NewSender(consumer, appLog)
		go func() {
			if err := app.Run(ctx); err != nil {
				log.Println(err)
				cancel()
			}
		}()

		<-ctx.Done()
	},
}
