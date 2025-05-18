package project_templates

const bootstrapFxTemplate = `package bootstrap

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"{{.Name}}/internal/config"
)

// Module provides core dependencies
var Module = fx.Options(
	fx.Provide(
		NewLogger,
		NewHTTPServer,
		config.GetConfig,
	),
)

func BuildApp() *fx.App {
	return fx.New(
		// Provide core dependencies
		Module,

		// Import application module
		fx.Provide(config.NewConfig),
		fx.Import("{{.Name}}/internal/app"),

		// Register lifecycle hooks
		fx.Invoke(RegisterHooks),
	)
}

func RegisterHooks(lc fx.Lifecycle, logger *zap.Logger) {
	// Register any global lifecycle hooks here
}
`

const bootstrapLoggerTemplate = `package bootstrap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"{{.Name}}/internal/config"
)

func NewLogger(cfg *config.Config) *zap.Logger {
	var zapConfig zap.Config

	if cfg.App.Debug {
		// Development logger configuration
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		// Production logger configuration
		zapConfig = zap.NewProductionConfig()
	}

	logger, _ := zapConfig.Build()
	return logger
}
`

const bootstrapHttpTemplate = `package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"{{.Name}}/internal/config"
)

func NewHTTPServer(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Lifecycle hooks
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
			logger.Info("Starting HTTP server", zap.String("addr", addr))
			
			go func() {
				if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
					logger.Error("Failed to start HTTP server", zap.Error(err))
				}
			}()
			
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping HTTP server")
			return e.Shutdown(ctx)
		},
	})

	return e
}
`

const bootstrapPostgresTemplate = `package bootstrap

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"{{.Name}}/internal/config"
)

func NewPostgresConnection(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host,
		cfg.Postgres.Port, cfg.Postgres.Database, cfg.Postgres.SSLMode)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to PostgreSQL",
				zap.String("host", cfg.Postgres.Host),
				zap.Int("port", cfg.Postgres.Port),
				zap.String("database", cfg.Postgres.Database))
			return pool.Ping(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing PostgreSQL connection")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func NewGoquDatabase() *goqu.Database {
	return goqu.New("postgres", nil)
}
`

const bootstrapRedisTemplate = `package bootstrap

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"{{.Name}}/internal/config"
)

func NewRedisClient(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to Redis",
				zap.String("host", cfg.Redis.Host),
				zap.Int("port", cfg.Redis.Port))
			return client.Ping(ctx).Err()
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Redis connection")
			return client.Close()
		},
	})

	return client, nil
}
`

const bootstrapKafkaTemplate = `package bootstrap

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"{{.Name}}/internal/config"
)

func NewKafkaWriter(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) *kafka.Writer {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Initializing Kafka writer",
				zap.Strings("brokers", cfg.Kafka.Brokers),
				zap.String("topic", cfg.Kafka.Topic))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing Kafka writer")
			return writer.Close()
		},
	})

	return writer
}

func NewKafkaReader(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.Topic,
		GroupID: cfg.Kafka.GroupID,
	})

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Initializing Kafka reader",
				zap.Strings("brokers", cfg.Kafka.Brokers),
				zap.String("topic", cfg.Kafka.Topic),
				zap.String("groupID", cfg.Kafka.GroupID))
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

const bootstrapGRPCTemplate = `package bootstrap

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"{{.Name}}/internal/config"
)

func NewGRPCServer(lc fx.Lifecycle, cfg *config.Config, logger *zap.Logger) *grpc.Server {
	server := grpc.NewServer()

	// Register services here
	// Example: pb.RegisterUserServiceServer(server, userService)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			logger.Info("Starting gRPC server", zap.String("addr", addr))

			go func() {
				if err := server.Serve(listener); err != nil {
					logger.Error("Failed to start gRPC server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping gRPC server")
			server.GracefulStop()
			return nil
		},
	})

	return server
}
`
