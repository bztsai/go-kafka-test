package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"sync"
)

func main() {
	topic := "hello"

	client, err := sarama.NewClient([]string{"localhost:19092"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Error(err)
		}
	}()

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Error(err)
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, p := range partitions {
		wg.Add(1)
		go func(p int32) {
			defer wg.Done()

			pc, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err := pc.Close(); err != nil {
					log.Error(err)
				}
			}()

			log.Infof("consuming %s:%d", topic, p)
			for msg := range pc.Messages() {
				fmt.Printf("%s:%d %s %s\n", topic, p, msg.Key, msg.Value)
			}
		}(p)
	}
	wg.Wait()
}
