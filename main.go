package main

import (
	"crypto/rand"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

func topic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}
	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	closeNotify := w.(http.CloseNotifier).CloseNotify()
	events := make(chan string)

	go runKafkaListener(vars["topic"], closeNotify, events)

	for {
		event, ok := <-events
		if !ok {
			return
		}
		fmt.Fprintf(w, "%s\n", event)
		flusher.Flush()
	}
}

func runKafkaListener(topic string, closeNotify <-chan bool, events chan string) {
	kafkaServers, exists := os.LookupEnv("KAFKA_SERVERS")
	if !exists {
		kafkaServers = "localhost:9092"
	}

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaServers,
		"group.id":          "kafka-web-" + pseudoUUID(),
		"auto.offset.reset": "latest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{topic}, nil)
	connectionClosed := false

	go func() {
		<-closeNotify
		close(events)
		c.Close()
		connectionClosed = true
	}()

	for {
		msg, err := c.ReadMessage(10 * time.Second)
		if err == nil {
			events <- string(msg.Value)
		} else {
			if connectionClosed {
				break
			}
		}
	}
}

func pseudoUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/topic/{topic}", topic)
	log.Println("Listen on port 8080...")
	http.ListenAndServe(":8080", r)
}
