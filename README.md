# Tests around Kafka and Go

Repository to experiment with different Kafka Go libraries.

## Notes

Update dependency:

```
$ go mod tidy
```

Topic admin:

```
$ kafka-topics --zookeeper localhost:22181 --list
$ kafka-topics --zookeeper localhost:22181 --create --topic hello --partitions 3 --replication-factor 2
$ kafka-topics --zookeeper localhost:22181 --delete --topic hello
```