package main

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var description = `Telnet client for test task

Example usage:
	go-telnet hostname port 
	go-telnet --timeout=10s hostname port 

Description params:
	hostname - string  name host for connect
	port - int port host for connect 

Description flags:
`
var ctx, cancel = context.WithCancel(context.Background())

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	timeout, host, port := parseFlagParams()
	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	SendToServer(ctx, client)
	ReceiveFromServer(ctx, client)

	go handlerSignal()
	<-ctx.Done()
}

func parseFlagParams() (timeout time.Duration, host string, port string) {
	pflag.DurationVar(&timeout, "timeout", time.Second*10, "timeout")
	pflag.Usage = func() {
		fmt.Fprint(os.Stderr, description)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	host = pflag.Arg(0)
	port = pflag.Arg(1)

	return
}

func handlerSignal() {
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGINT)
	for {
		select {
		case <-signalCh:
			cancel()
			return
		case <-ctx.Done():
			return
		}
	}
}
