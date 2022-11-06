package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/pflag"
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

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	var ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()
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
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			if err := client.Send(); err != nil {
				log.Println(err)
				stop()
				return
			}
		}
	}()
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			if err := client.Receive(); err != nil {
				log.Println(err)
				stop()
				return
			}
		}
	}()

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
