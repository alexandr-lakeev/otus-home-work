package internalampq

import (
	"context"
	"encoding/json"

	appsender "github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/app/sender"
	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Consumer struct {
	config  config.AmpqConf
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewConsumer(config config.AmpqConf) *Consumer {
	return &Consumer{
		config: config,
	}
}

func (c *Consumer) Connect(ctx context.Context) (err error) {
	c.conn, err = amqp.Dial(c.config.URL)
	if err != nil {
		return err
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return err
	}

	return c.init()
}

func (c *Consumer) Close(ctx context.Context) error {
	return c.conn.Close()
}

func (c *Consumer) Consume(ctx context.Context, fn appsender.Worker, threads int) error {
	msgCh, err := c.channel.Consume(
		c.config.QueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	ch := make(chan appsender.Event)

	for i := 0; i < threads; i++ {
		go fn(ctx, ch)
	}

	for {
		select {
		case msg := <-msgCh:
			var event appsender.Event

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				return err
			}

			ch <- event
			msg.Ack(true)
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Consumer) init() error {
	_, err := c.channel.QueueDeclare(
		c.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return c.channel.QueueBind(c.config.QueueName, "", c.config.ExchangeName, false, nil)
}
