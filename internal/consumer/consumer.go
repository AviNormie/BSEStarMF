package consumer

import (
	"context"
	
	"time"

	"github.com/segmentio/kafka-go"
	"sapphirebroking.com/sapphire_mf/internal/config"
	"sapphirebroking.com/sapphire_mf/internal/util"
)

type Consumer struct {
	reader       *kafka.Reader
	config       *config.KafkaConfig
	handler      MessageHandler
	logger       util.Logger
	shutdownChan chan struct{}
}

func NewConsumer(cfg *config.Config, handler MessageHandler, logger util.Logger) (*Consumer, error) {
	brokers := cfg.GetKafkaBrokers() // Get effective brokers
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     cfg.Kafka.GroupID,
		Topic:       cfg.Kafka.Topic,
		StartOffset: kafka.FirstOffset, // or kafka.LastOffset based on cfg.AutoOffsetReset
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
	})

	return &Consumer{
		reader:       reader,
		config:       cfg.Kafka,
		handler:      handler,
		logger:       logger,
		shutdownChan: make(chan struct{}),
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	defer c.reader.Close()
	c.logger.Info("Starting Kafka consumer on topic: %s", c.config.Topic)

	for {
		select {
		case <-c.shutdownChan:
			c.logger.Info("Shutting down Kafka consumer on topic: %s", c.config.Topic)
			return
		case <-ctx.Done():
			c.logger.Info("Context cancelled, Shutting down Kafka consumer on topic: %s", c.config.Topic)
			return
		default:
			// Set a timeout for reading messages
			msgCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
			msg, err := c.reader.ReadMessage(msgCtx)
			cancel()

			if err != nil {
				if err == context.DeadlineExceeded {
					// Timeout is normal, continue polling
					continue
				}
				c.logger.Error("Error reading message: %v", err)
				continue
			}

			c.logger.Info("Received message from topic %s [%d] at offset %d: %s", 
				msg.Topic, msg.Partition, msg.Offset, string(msg.Value))

			handleCtx, handleCancel := context.WithTimeout(ctx, 5*time.Second)
			err = c.handler.HandleMessage(handleCtx, msg.Key, msg.Value)
			handleCancel()

			if err != nil {
				c.logger.Error("Error processing message: %v. Message: %s", err, string(msg.Value))
			} else {
				// With segmentio/kafka-go, commits are handled automatically by the reader
				c.logger.Info("Successfully processed message at offset %d", msg.Offset)
			}
		}
	}
}

func (c *Consumer) Shutdown() {
	close(c.shutdownChan)
}