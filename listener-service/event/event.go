package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange declares a durable topic exchange named "logs_topic"
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

// declareRandomQueue declares an exclusive, auto-delete queue with a random name
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name - empty string means a random name is generated
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
}


