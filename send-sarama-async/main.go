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
	defer func() {
		if err := client.Close(); err != nil {
			log.Error(err)
		}
	}()

	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer func() {
			if err := producer.Close(); err != nil {
				log.Error(err)
			}
		}()

		for i := 0; i < 10; i++ {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Key:   sarama.StringEncoder(fmt.Sprintf("async-sarama-%d", i)),
				Value: sarama.StringEncoder("hi"),
			}
			producer.Input() <- msg
		}
	}()

	successesCh := producer.Successes()
	errorsCh := producer.Errors()
	for successesCh != nil && errorsCh != nil {
		select {
		case msg, ok := <-successesCh:
			if ok {
				log.Infof("%s success", msg.Key)
			} else {
				successesCh = nil
			}
		case msg, ok := <-errorsCh:
			if ok {
				log.Infof("%s error - %v", msg.Msg.Key, msg.Err)
			} else {
				errorsCh = nil
			}
		}
	}
}
