package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func main() {
	topic := "hello"

	conf := sarama.NewConfig()
	conf.Producer.Compression = sarama.CompressionGZIP
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true
	conf.Producer.RequiredAcks = sarama.WaitForAll

	client, err := sarama.NewClient([]string{"localhost:19092"}, conf)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer producer.Close()

		for i := 0; i < 10; i++ {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key: sarama.StringEncoder(fmt.Sprintf("async-sarama-%d", i)),
				Value: sarama.StringEncoder("hi"),
			}
			producer.Input() <- msg
		}
	}()

	successesDone := false
	errorsDone := false
	for !successesDone && !errorsDone {
		select {
		case msg, ok := <-producer.Successes():
			if ok {
				log.Infof("%s success", msg.Key)
			} else {
				successesDone = true
			}
		case msg, ok := <-producer.Errors():
			if ok {
				log.Infof("%s error - %v", msg.Msg.Key, msg.Err)
			} else {
				errorsDone = true
			}
		}
	}
}