package config

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shordem/api.thryvo/lib/constants"
)

// EventInterface is an interface for the Event struct
type EventInterface interface {
	Publish(ctx context.Context, exchange, key string, msg []byte) error
	Consume(ctx context.Context, exchange, key string) (<-chan amqp.Delivery, error)
}

// event is a struct that holds the connection and channel to the RabbitMQ server
type event struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewEvent creates a new Event struct
func NewEvent(env constants.Env) EventInterface {
	event := &event{}
	event.init(env)
	return event
}

// init initializes the connection and channel to the RabbitMQ server
func (e *event) init(env constants.Env) {
	conn, err := amqp.Dial(env.RABBITMQ_SERVER)
	e.FailOnError(err, "Failed to connect")

	ch, err := conn.Channel()
	e.FailOnError(err, "Failed to open a channel")

	defer ch.Close()

	log.Println("Connected to RabbitMQ")

	e.conn = conn
	e.channel = ch
}

// Publish sends a message to the RabbitMQ server
func (e *event) Publish(ctx context.Context, exchange, key string, msg []byte) error {
	err := e.channel.Publish(exchange, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        msg,
	})
	if err != nil {
		return err
	}

	return nil
}

// Consume receives a message from the RabbitMQ server
func (e *event) Consume(ctx context.Context, exchange, key string) (<-chan amqp.Delivery, error) {
	msgs, err := e.channel.Consume(
		key,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (e *event) DeclareQueue(name string) (amqp.Queue, error) {
	q, err := e.channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)

	e.FailOnError(err, "Failed to declare a queue")

	return q, nil
}

// FailOnError logs a fatal error if the given error is not nil
func (e *event) FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("RabbitMQ: %s: %s", msg, err)
	}
}
