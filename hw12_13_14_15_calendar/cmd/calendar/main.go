package main

import (
	"log"

	"github.com/AlekseyPorandaykin/otus_hw_go/hw12_13_14_15_calendar/cmd/calendar/cmd"
	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
