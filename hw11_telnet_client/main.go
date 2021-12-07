package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	telent := NewTelnetClient("127.0.0.1:4242", 120*time.Second, os.Stdin, os.Stdout)
	if err := telent.Connect(); err != nil {
		fmt.Println(err)
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
