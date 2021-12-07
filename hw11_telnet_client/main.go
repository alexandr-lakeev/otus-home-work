package main

import (
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Timeout string `long:"timeout" description:"Client timeout"`
}

func main() {
	args, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	address := net.JoinHostPort(args[0], args[1])
	timeout, err := time.ParseDuration(opts.Timeout)
	if err != nil {
		log.Fatal(err)
	}

	telent := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := telent.Connect(); err != nil {
		log.Fatal(err)
	}
	defer telent.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		telent.Receive()
	}()

	go func() {
		defer wg.Done()
		telent.Send()
	}()

	wg.Wait()
}
