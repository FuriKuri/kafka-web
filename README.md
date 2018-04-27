# kafka-web
Simple web server to view into kafka topics

## Usage

### /topic/{topic-name} - Listen to topic

Use `http://localhost:8080/topic/{topic-name}` to get all new message from a topic `topic-name`

## Deploy

```
$ docker run -p 8080:8080 -e KAFKA_SERVERS=192.168.0.122:9092 furikuri/kafka-web
```

Open a browser or use curl to listen to the topic:
```
$ curl 192.168.0.122:8080/topic/hello-world
This is a message
Hi mom
```

You can use [kafka-cat](https://github.com/edenhill/kafkacat) to produce messages:
```
$ kafkacat -P -b 192.168.0.122:9092 -t hello-world
```