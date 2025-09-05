package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
)

var latestPrice []byte

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/crypto/latest", func(w http.ResponseWriter, r *http.Request) {
		if latestPrice == nil {
			http.Error(w, "No data available", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(latestPrice)
	}).Methods("GET")

	go startConsumer()

	log.Println("🚀 Consumer API running at http://localhost:8081/crypto/latest")
	log.Fatal(http.ListenAndServe(":8081", r))
}

func startConsumer() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "crypto_prices",
		GroupID: "crypto_group",
	})
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.TODO())
		if err != nil {
			log.Println("Error reading from Kafka:", err)
			continue
		}
		latestPrice = msg.Value
		log.Println("✅ Consumed:", string(msg.Value))
	}
}
