package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	amgp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn  *amgp.Connection
	queue string
}

func NewConsumer(conn *amgp.Connection) (Consumer, error) {
	c := Consumer{
		conn: conn,
	}

	// consumer.setup() is for creating the queue and binding it to the exchange

	err := c.setup()

	if err != nil {
		return Consumer{}, err
	}

	return c, nil
}

func (c *Consumer) setup() error {
	// channel, err := c.conn.Channel() is for creating a channel (we'll use it to create the queue)
	channel, err := c.conn.Channel()

	if err != nil {
		return err
	}

	// declareExchange(channel) is for creating the exchange

	return declareExchange(channel)
}

type Payload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Active    int    `json:"active"`
	Type      string `json:"type"`
}

func (c *Consumer) Listen(topics []string) error {
	fmt.Println("Listening for messages... consumer.go")
	// channel, err := c.conn.Channel() is for creating a channel (we'll use it to create the queue)
	channel, err := c.conn.Channel()

	if err != nil {
		return err
	}

	// defer channel.Close() will not work here, because the program will exit before it is called
	defer channel.Close()

	//declareRandomQueue(channel) is for creating the queue (the name of the queue will be a random string)

	queue, err := declareRandomQueue(channel)

	if err != nil {
		return err
	}

	// for _, topic := range topics is for looping through the topics we want to listen to (e.g. "post.created", "post.deleted", "post.updated")
	fmt.Println("topics: ", topics)
	for _, s := range topics {
		fmt.Println("s: ", s)
		fmt.Println("queue.Name: ", queue.Name)
		// bindQueue(channel, queue, topic) is for binding the queue to the exchange (so that the queue will receive messages published to the exchange)

		// channel.QueueBind(queue, s, "events", false, nil) is for binding the queue to the exchange (so that the queue will receive messages published to the exchange)
		channel.QueueBind(queue.Name, s, "events_user", false, nil)

		if err != nil {
			return err
		}

	}

	// msgs, err := channel.Consume(queue, "", true, false, false, false, nil) is for consuming messages from the queue (we'll use it for receiving messages)
	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	fmt.Println("messages: ", messages)

	if err != nil {
		return err
	}

	// forever:= make(chan bool) is for blocking the main thread (so that the program doesn't exit before we receive any messages)

	forever := make(chan bool)

	// go func() is for running a function in a goroutine (a lightweight thread)
	go func() {
		for d := range messages {
			var payload Payload

			fmt.Println("d.Body: ", d.Body)
			_ = json.Unmarshal(d.Body, &payload)

			fmt.Println("payload: ", payload)

			go handlePayload(payload)

		}
	}()

	fmt.Println("Listening for messages...")
	// <-forever is for blocking the main thread (so that the program doesn't exit before we receive any messages)
	<-forever

	return nil
}

// func handlePayload(payload Payload) is for handling the payload (printing it to the console)

func handlePayload(payload Payload) {

	switch payload.Type {
	case "user.created":
		if err := userCreated(payload); err != nil {
			fmt.Printf("error creating post: %s", err)
		}

	case "user.deleted":
		fmt.Println("post deleted")
	case "user.updated":
		fmt.Println("post updated")
	default:
		err := fmt.Errorf("unknown event: %s", payload.Type)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// func logEvent(payload Payload) is for logging the event to the database (using the event-service)

type userCreatedPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Active    int    `json:"active"`
}

func userCreated(payload Payload) error {

	payloadPost := userCreatedPayload{
		Email:     payload.Email,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Password:  payload.Password,
		Active:    payload.Active,
	}

	jsonDat, _ := json.MarshalIndent(payloadPost, "", "\t")

	request, err := http.NewRequest("POST", "http://user-service/v1/user", bytes.NewBuffer(jsonDat))

	if err != nil {

		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {

		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {

		return errors.New("unexpected status from post-service: " + response.Status)
	}

	return nil
}
