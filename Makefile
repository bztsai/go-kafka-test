.PHONY: ps up down topics-list topics-create-hello compile clean

DOCKER_COMPOSE=docker-compose -f ./docker/kafka-cluster/docker-compose.yml
KAFKA_TOPICS=${DOCKER_COMPOSE} exec kafka-1 kafka-topics --zookeeper localhost:22181

ps:
	${DOCKER_COMPOSE} ps

up:
	${DOCKER_COMPOSE} up -d

down:
	${DOCKER_COMPOSE} down

topics-list:
	${KAFKA_TOPICS} --list

topics-create-hello:
	${KAFKA_TOPICS} --create --topic hello --partitions 3 --replication-factor 2

compile: bin/kafkacat bin/send-sarama-async

clean:
	rm -rf bin

bin/kafkacat:
	cd kafkacat && go build -o ../bin/kafkacat ./...

bin/send-sarama-async:
	cd send-sarama-async && go build -o ../bin/send-sarama-async ./...