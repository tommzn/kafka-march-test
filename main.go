package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func main() {

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	bootstrapServers := flag.String("bootstrap-servers", "localhost", "Pass Kafka server(s)")
	topic := flag.String("topic", "hello-world", "A topic you want to listen to")
	flag.Parse()

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": *bootstrapServers,
		"group.id":          "kafka-march-test",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}
	c.SubscribeTopics([]string{*topic}, nil)

	run := true
	go func() {
		<-sigs
		run = false
	}()

	for run {
		msg, err := c.ReadMessage(time.Second)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}
	c.Close()
}
