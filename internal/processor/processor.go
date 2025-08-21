package processor

import (
	"context"
	"encoding/json"
	"fmt"

	"sapphirebroking.com/sapphire_mf/internal/util"
)

type Processor struct {
	logger util.Logger
}

func NewProcessor(logger util.Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

func (p *Processor) HandleMessage(ctx context.Context, key, value []byte) error {
	p.logger.Info("Processing MF message with key: %s", string(key))

	// Basic JSON parsing to check if message is valid
	var message map[string]interface{}
	if err := json.Unmarshal(value, &message); err != nil {
		p.logger.Error("Failed to unmarshal MF message: %v", err)
		return fmt.Errorf("failed to unmarshal MF message: %w", err)
	}

	p.logger.Info("Successfully processed MF message: %+v", message)
	return nil
}