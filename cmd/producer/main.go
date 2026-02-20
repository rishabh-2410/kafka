package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

func main() {

	// Load env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a producer instance
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		// User specific properties to set
		"bootstrap.servers": os.Getenv("BROKER"),
		"sasl.username":     os.Getenv("KEY"),
		"sasl.password":     os.Getenv("SECRET"),

		// Fixed properties
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"acks":              "all",
	})

	if err != nil {
		fmt.Printf("Failed to create the producer: %s", err.Error())
		os.Exit(1)
	}

	// Go routine to handle message delivery reports and other event types(errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition.Error.Error())
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value =%s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	users := [...]string{"eabara", "jsmith", "sgarcia", "jbernard", "htanaka", "awalther"}
	items := [...]string{"book", "alarm clock", "t-shirts", "gift card", "batteries"}
	topic := "purchases"

	// Produce mock data
	for range 10 {
		key := users[rand.Intn(len(users))]  // Select a random user as key
		data := items[rand.Intn(len(items))] // Select a random item purchased as data

		// Produce the actuall message and send to a internal queue (buffer)
		// A background goroutine inside kafka producer runs continuously to send the
		// messages from the buffer to kafka
		// App -> Buffer -> Kafka

		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(key),
			Value: []byte(data),
		}, nil)
	}

	// Time to wait for the above background goroutine inside kafka producer
	// to finish sending all messages to the topic
	p.Flush(15 * 1000)

	// Close the producer instance
	// The Producer object or its channels are no longer usable after this call
	p.Close()
}
