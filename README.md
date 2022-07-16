kafka-go-demo
-----

## Docker container

Run kafka container:
```
docker-compose up -d

docker container ls
```

Kafka commands location:
```
docker exec -it 3d35558a778e /bin/bash

/opt/bitnami/kafka/bin
```

## Build and Run
```
go mod tidy

go mod vendor

go run ./segmentio/main.go
```

## Kafka CLI
```
docker exec -it 5a5abe94a0c3 /opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list

docker exec -it 5a5abe94a0c3 /opt/bitnami/kafka/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic my_topic

```


## Resources

### MSK
https://aws.amazon.com/msk/what-is-kafka/
https://docs.aws.amazon.com/msk/latest/developerguide/what-is-msk.html


### Kafka container
https://hub.docker.com/r/bitnami/kafka


### Confluent
https://github.com/confluentinc/confluent-kafka-go
https://reposhub.com/go/miscellaneous/confluentinc-confluent-kafka-go.html
https://github.com/confluentinc/confluent-kafka-go/blob/master/examples/stats_example/stats_example.go


### Segmentio
https://github.com/segmentio/kafka-go
