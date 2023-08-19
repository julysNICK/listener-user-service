package main

import (
	"fmt"
	"listener-user-service/event"
	"log"
	"math"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	rabbitConn, err := connect()

	if err != nil {
		log.Fatalf("could not connect to rabbitmq: %v", err)
		os.Exit(1)
	}
	// defer rabbitConn.Close() will not work here, because the program will exit before it is called
	defer rabbitConn.Close()

	log.Println("Listening for messages... User")

	// channel is for sending and receiving messages from the queue (we'll only use it for receiving messages)

	//consume is a channel that will deliver messages from the queue to us (we'll only use it for receiving messages)

	consumer, err := event.NewConsumer(rabbitConn)

	if err != nil {
		panic(err)
	}

	// consumer.Listen([]string{"user.created", "user.deleted", "user.updated"}) is for listening to messages published to the exchange (we'll use it for receiving messages)
	err = consumer.Listen([]string{"user.created", "user.deleted", "user.updated"})
	if err != nil {
		panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	// var counts int64 is for counting the number of messages we've received from the queue
	var counts int64
	// var backOff is for exponential backoff (we'll use it to reconnect to RabbitMQ if the connection is lost)
	var backOff = 1 * time.Second

	//var conn *amqp.Connection is for storing the connection to RabbitMQ (we'll use it to create a channel)
	var connection *amqp.Connection
	//for loop will keep trying to connect to RabbitMQ until it succeeds (or until the program is terminated)
	for {
		// c, err := amqp.Dial("amqp://guest:guest@localhost:5672/") is for connecting to RabbitMQ
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")

		if err != nil {
			log.Println("could not connect to rabbitmq, retrying... user")
			counts++
		} else {
			log.Println("connected to rabbitmq")
			connection = c
			break
		}

		if counts > 10 {
			fmt.Println("could not connect to rabbitmq after 10 tries, exiting... user")
			return nil, err
		}

		backOff = time.Duration(math.Pow(2, float64(counts))) * time.Second

		log.Printf("waiting %v seconds before trying again user", backOff.Seconds())

		time.Sleep(backOff)

		continue
	}
	return connection, nil
}
