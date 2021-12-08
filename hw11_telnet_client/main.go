package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	timeout := flag.Duration("timeout", 0, "Client timeout")
	flag.Parse()

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))
	telent := NewTelnetClient(address, *timeout, os.Stdin, os.Stdout)

	if err := telent.Connect(); err != nil {
		log.Fatal(err)
	}
	defer telent.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		defer cancel()
		telent.Receive()
	}()

	go func() {
		defer cancel()
		telent.Send()
	}()

	<-ctx.Done()
}
