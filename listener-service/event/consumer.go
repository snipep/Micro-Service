package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer represents an AMQP consumer
type Consumer struct {
	conn      	*amqp.Connection
	queueName 	string
}

// NewConsumer returns a new Consumer
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer {
		conn: conn,
	}
	err:= consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}

// setup sets up the exchange and queue on the AMQP broker
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	// declare the exchange
	return declareExchange(channel)
}

// Payload represents a message payload
type Payload struct {
	Name 	string		`json:"name"`
	Data 	string		`json:"data"`
}

// Listen consumes messages from the broker and handles them
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil{
		return err
	}
	defer ch.Close()

	// declare a random queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// bind the queue to the exchange for each topic
	for _, s := range topics {
		err = ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	// consume messages from the queue
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// start consuming messages
	forever := make(chan bool)
	go func(){
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go consumer.handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message on exchange [Exchange, Qeueue] [logs_topics, %s]\n", q.Name)
	<-forever

	return nil
}

// handlePayload handles a message payload
func (consumer *Consumer) handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	case "auth":
		// authenticate
	
	default:
		err := logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}	
	}
	panic("unimplemenlted")
}

// logEvent logs an event to the logger service
func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// set the content type to application/json
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// check the response status code
	if response.StatusCode != http.StatusAccepted {
		return err
	}
	return nil 
}
