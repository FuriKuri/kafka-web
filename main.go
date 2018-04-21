package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func topic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Printf("URL %s", vars["topic"])
	flusher, ok := w.(http.Flusher)
	if !ok {
		panic("expected http.ResponseWriter to be an http.Flusher")
	}

	w.Header().Set("Connection", "Keep-Alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	for i := 1; i <= 10; i++ {
		fmt.Fprintf(w, "Chunk #%d\n", i)
		flusher.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/topic/{topic}", topic)
	log.Println("Listen on port 8080...")
	http.ListenAndServe(":8080", r)
}
