package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/realPointer/orders/internal/model"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	order := model.Order{
		OrderUID:    "test123",
		TrackNumber: "TEST123",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:  "Test Testov",
			Phone: "+9720000000",
			Email: "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction: "b563feb7b2b84b6test",
			Currency:    "USD",
			Provider:    "wbpay",
			Amount:      1817,
			PaymentDt:   time.Now(),
		},
		Items: []model.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
			},
		},
		DateCreated: time.Now(),
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Fatalf("Failed to marshal order: %s", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: "orders",
		Value: sarama.StringEncoder(orderJSON),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		log.Fatalf("Failed to send message: %s", err)
	}

	log.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
