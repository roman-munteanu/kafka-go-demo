kafka-go-demo
-----

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

Build and run
```
go mod tidy

go mod vendor

go run ./segmentio/main.go
```

Kafka CLI
```
docker exec -it 3d35558a778e /opt/bitnami/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list

docker exec -it 3d35558a778e /opt/bitnami/kafka/bin/kafka-console-producer.sh --broker-list localhost:9092 --topic my_topic

```