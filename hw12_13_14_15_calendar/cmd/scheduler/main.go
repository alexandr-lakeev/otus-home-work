package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/infrastructure/storage"
	"github.com/streadway/amqp"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler.dev.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewSchedulerConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(config)

	storage, err := storage.New(config.Storage)
	if err != nil {
		log.Fatalln(err)
	}

	scheduler := scheduler.New(storage)

	ctx := context.Background()
	scheduler.Notify(ctx, 5*time.Minute)

	// ticker := time.NewTicker(5 * time.Second)

	// for {
	// 	log.Println("tick")

	// 	publish()

	// 	<-ticker.C
	// }
}

func publish() {
	connection, err := amqp.Dial("amqp://user:password@rabbit:5672/")
	if err != nil {
		fmt.Printf("Dial: %s", err)
		return
	}
	defer connection.Close()

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		fmt.Printf("Channel: %s", err)
		return
	}

	if err := channel.ExchangeDeclare(
		"published_events", // name
		"fanout",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // noWait
		nil,                // arguments
	); err != nil {
		fmt.Printf("Exchange Declare: %s", err)
		return
	}

	// TODO move to consumer?
	_, err = channel.QueueDeclare(
		"published_events.sender",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("Queue: %s", err)
		return
	}

	if err := channel.QueueBind("published_events.sender", "", "published_events", false, nil); err != nil {
		fmt.Printf("QueueBind: %s", err)
		return
	}

	if err = channel.Publish(
		"published_events",
		"",
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte("Foo"),
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	); err != nil {
		fmt.Printf("Exchange Publish: %s", err)
		return
	}
}
