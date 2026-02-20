package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error reading env variables")
	}

	// Create a new consumer instance
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		// User specific properties to set
		"bootstrap.servers": os.Getenv("BROKER"),
		"sasl.username":     os.Getenv("KEY"),
		"sasl.password":     os.Getenv("SECRET"),

		// Fixed properties
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"auto.offset.reset": "earliest",
		"group.id":          "kafka-consumer-groud-0",
	})

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err.Error())
		os.Exit(1)
	}

	topic := "purchases"
	err = c.SubscribeTopics([]string{topic}, nil)

	// Setup a channel for handling Ctrl+C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case signal := <-sigchan:
			fmt.Printf("Caught signal %v, terminating...\n", signal)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by customer
				continue
			}

			fmt.Printf("Consumed event from topic: %s, key = %-10s, value = %-10s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value),
			)
		}
	}

	// Close the consumer instance
	c.Close()
}
