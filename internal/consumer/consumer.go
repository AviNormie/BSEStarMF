package consumer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"sapphirebroking.com/sapphire_mf/internal/config"
	"sapphirebroking.com/sapphire_mf/internal/util"
)

type Consumer struct {
	consumer     *kafka.Consumer
	config       *config.KafkaConfig
	handler      MessageHandler
	logger       util.Logger
	shutdownChan chan struct{}
}

func NewConsumer(cfg *config.KafkaConfig, handler MessageHandler, logger util.Logger) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     strings.Join(cfg.Brokers, ","),
		"group.id":              cfg.GroupID,
		"auto.offset.reset":     cfg.AutoOffsetReset,
		"enable.auto.commit":    false,
		"session.timeout.ms":    6000,
		"heartbeat.interval.ms": 2000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	logger.Info("Kafka consumer created for group: %s, topic: %s", cfg.GroupID, cfg.Topic)

	return &Consumer{
		consumer:     c,
		config:       cfg,
		handler:      handler,
		logger:       logger,
		shutdownChan: make(chan struct{}),
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	defer c.consumer.Close()
	c.logger.Info("Starting Kafka consumer on topic: %s", c.config.Topic)

	err := c.consumer.SubscribeTopics([]string{c.config.Topic}, nil)
	if err != nil {
		c.logger.Fatal("Failed to subscribe to Kafka topic %s: %v", c.config.Topic, err)
	}

	for {
		select {
		case <-c.shutdownChan:
			c.logger.Info("Shutting down Kafka consumer on topic: %s", c.config.Topic)
			return
		case <-ctx.Done():
			c.logger.Info("Context cancelled, Shutting down Kafka consumer on topic: %s", c.config.Topic)
			return
		default:
			ev := c.consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				c.logger.Info("Received message from topic %s [%d] at offset %v: %s", *e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset, string(e.Value))
				handleCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
				err := c.handler.HandleMessage(handleCtx, e.Key, e.Value)
				cancel()

				if err != nil {
					c.logger.Error("Error processing message: %v. Message: %s", err, string(e.Value))
				} else {
					_, err := c.consumer.CommitMessage(e)
					if err != nil {
						c.logger.Error("Failed to commit offset: %v", err)
					} else {
						c.logger.Info("Committed offset for message at offset %v", e.TopicPartition.Offset)
					}
				}
			case kafka.Error:
				if e.IsFatal() {
					c.logger.Fatal("Kafka consumer encountered fatal error: %v", e)
				} else {
					c.logger.Error("Kafka consumer encountered an error: %v", e)
				}
			}
		}
	}
}

func (c *Consumer) Shutdown() {
	close(c.shutdownChan)
}