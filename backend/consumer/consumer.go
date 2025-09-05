package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/segmentio/kafka-go"
)

type CryptoPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}

var latestPrice []byte

func main() {
	// Start Kafka consumer
	go startConsumer()

	// Serve React static files
	frontendPath := filepath.Join("frontend", "build")
	if _, err := os.Stat(frontendPath); os.IsNotExist(err) {
		log.Println("⚠️ React build folder not found. Create it later with React code.")
	} else {
		fs := http.FileServer(http.Dir(frontendPath))
		http.Handle("/", fs)
		log.Println("✅ Serving React UI from:", frontendPath)
	}

	// API endpoint for latest crypto price
	http.HandleFunc("/crypto/latest", func(w http.ResponseWriter, r *http.Request) {
		if latestPrice == nil {
			http.Error(w, "No data available", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(latestPrice)
	})

	log.Println("🚀 Consumer API + React UI running at http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
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
			log.Println("❌ Error reading from Kafka:", err)
			continue
		}
		latestPrice = msg.Value
		log.Println("✅ Consumed:", string(msg.Value))
	}
}
