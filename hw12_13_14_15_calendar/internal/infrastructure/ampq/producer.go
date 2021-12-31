package internalampq

import (
	"context"
	"encoding/json"

	"github.com/alexandr-lakeev/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/streadway/amqp"
)

type Producer struct {
	config  config.AmpqConf
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer(config config.AmpqConf) *Producer {
	return &Producer{
		config: config,
	}
}

func (p *Producer) Connect(ctx context.Context) (err error) {
	p.conn, err = amqp.Dial(p.config.URL)
	if err != nil {
		return err
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return err
	}

	return p.init()
}

func (p *Producer) Close(ctx context.Context) error {
	return p.conn.Close()
}

func (p *Producer) Produce(ctx context.Context, data interface{}) error {
	channel, err := p.conn.Channel()
	if err != nil {
		return err
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return channel.Publish(
		p.config.ExchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "application/json",
			ContentEncoding: "utf-8",
			Body:            body,
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		},
	)
}

func (p *Producer) init() error {
	if err := p.channel.ExchangeDeclare(
		p.config.ExchangeName,
		p.config.ExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	); err != nil {
		return err
	}

	// TODO move to sender consumer?
	_, err := p.channel.QueueDeclare(
		p.config.QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// TODO move to sender consumer?
	return p.channel.QueueBind(p.config.QueueName, "", p.config.ExchangeName, false, nil)
}
