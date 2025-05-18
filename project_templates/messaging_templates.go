package project_templates

const kafkaTemplate = `package kafka

import (
	"context"
	
	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	
	"{{.Name}}/internal/config"
)


type KafkaConfig struct {
	Brokers []string
	Topic   string
}


func NewKafkaConfig(cfg *config.Config) *KafkaConfig {
	return &KafkaConfig{
		Brokers: []string{cfg.GetEnv("KAFKA_BROKER", "localhost:9092")},
		Topic:   cfg.GetEnv("KAFKA_TOPIC", "users"),
	}
}


func NewKafkaWriter(lc fx.Lifecycle, cfg *KafkaConfig, logger *zap.Logger) *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
	})
	
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Initializing Kafka writer", 
				zap.Strings("brokers", cfg.Brokers), 
				zap.String("topic", cfg.Topic))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Kafka writer")
			return writer.Close()
		},
	})
	
	return writer
}


func NewKafkaReader(lc fx.Lifecycle, cfg *KafkaConfig, logger *zap.Logger) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: "app-consumer",
	})
	
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Initializing Kafka reader", 
				zap.Strings("brokers", cfg.Brokers), 
				zap.String("topic", cfg.Topic))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Kafka reader")
			return reader.Close()
		},
	})
	
	return reader
}
`

const kafkaUserEventsTemplate = `package kafka

import (
	"context"
	"encoding/json"
	
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	
	"{{.Name}}/internal/domain"
)


type UserEventType string

const (
	UserCreated UserEventType = "user_created"
	UserUpdated UserEventType = "user_updated"
	UserDeleted UserEventType = "user_deleted"
)


type UserEvent struct {
	Type UserEventType ` + "`json:\"type\"`" + `
	User *domain.User  ` + "`json:\"user\"`" + `
}


type UserEventPublisher struct {
	writer *kafka.Writer
	logger *zap.Logger
}


func NewUserEventPublisher(writer *kafka.Writer, logger *zap.Logger) *UserEventPublisher {
	return &UserEventPublisher{
		writer: writer,
		logger: logger,
	}
}


func (p *UserEventPublisher) Publish(ctx context.Context, eventType UserEventType, user *domain.User) error {
	event := UserEvent{
		Type: eventType,
		User: user,
	}
	
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	
	err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(user.ID),
		Value: data,
	})
	
	if err != nil {
		p.logger.Error("Failed to publish user event", 
			zap.String("type", string(eventType)),
			zap.String("user_id", user.ID),
			zap.Error(err))
		return err
	}
	
	p.logger.Info("Published user event", 
		zap.String("type", string(eventType)),
		zap.String("user_id", user.ID))
	
	return nil
}


type UserEventConsumer struct {
	reader *kafka.Reader
	logger *zap.Logger
}


func NewUserEventConsumer(reader *kafka.Reader, logger *zap.Logger) *UserEventConsumer {
	return &UserEventConsumer{
		reader: reader,
		logger: logger,
	}
}


func (c *UserEventConsumer) Start(ctx context.Context) {
	go func() {
		c.logger.Info("Starting user event consumer")
		
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					c.logger.Error("Failed to read message", zap.Error(err))
					continue
				}
				
				var event UserEvent
				if err := json.Unmarshal(msg.Value, &event); err != nil {
					c.logger.Error("Failed to unmarshal user event", zap.Error(err))
					continue
				}
				
				c.logger.Info("Received user event", 
					zap.String("type", string(event.Type)),
					zap.String("user_id", event.User.ID))
				
				
				switch event.Type {
				case UserCreated:
					
				case UserUpdated:
					
				case UserDeleted:
					
				}
			}
		}
	}()
}
`
