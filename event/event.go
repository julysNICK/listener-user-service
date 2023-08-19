package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// func declareExchange(channel *amqp.Channel) error is for creating the exchange (if it doesn't already exist)
func declareExchange(channel *amqp.Channel) error {
	// channel.ExchangeDeclare("events", "topic", true, false, false, false, nil) is for creating the exchange
	return channel.ExchangeDeclare("events_user", "topic", true, false, false, false, nil)
}

// func declareRandomQueue(channel *amqp.Channel) (string, error) is for creating the queue (the name of the queue will be a random string) and binding it to the exchange
func declareRandomQueue(channel *amqp.Channel) (amqp.Queue, error) {
	return channel.QueueDeclare("", false, false, true, false, nil)

}
