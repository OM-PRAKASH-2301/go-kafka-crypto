package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

type CryptoPrice struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Time   string  `json:"time"`
}

func main() {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "crypto_prices",
	})
	defer writer.Close()

	rand.Seed(time.Now().UnixNano())

	for {
		price := CryptoPrice{
			Symbol: "BTCUSDT",
			Price:  100000 + rand.Float64()*10000,
			Time:   time.Now().Format(time.RFC3339),
		}

		data, _ := json.Marshal(price)
		err := writer.WriteMessages(context.TODO(),
			kafka.Message{
				Key:   []byte(price.Symbol),
				Value: data,
			},
		)
		if err != nil {
			log.Println("❌ Error writing to Kafka:", err)
		} else {
			log.Println("✅ Produced:", string(data))
		}

		time.Sleep(2 * time.Nanosecond)
	}
}
