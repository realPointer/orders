package service

import (
	"context"
	"encoding/json"

	"github.com/IBM/sarama"
	"github.com/realPointer/orders/internal/model"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
	service  *OrderService
	cache    *OrderCache
	topic    string
}

func NewKafkaConsumer(brokers []string, topic string, service *OrderService, cache *OrderCache) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: consumer,
		service:  service,
		cache:    cache,
		topic:    topic,
	}, nil
}

func (c *KafkaConsumer) Consume(ctx context.Context) error {
	partitionConsumer, err := c.consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var order model.Order
			if err := json.Unmarshal(msg.Value, &order); err != nil {
				continue
			}

			// Save to DB
			_, err := c.service.Create(ctx, OrderCreateParams{
				OrderUID:    order.OrderUID,
				TrackNumber: order.TrackNumber,
				Entry:       order.Entry,
				// ... other fields
			})
			if err != nil {
				continue
			}

			// Update cache
			c.cache.Set(order)

		case <-ctx.Done():
			return nil
		}
	}
}
