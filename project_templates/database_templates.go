package project_templates

const postgresTemplate = `package postgres

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


type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}


func NewPostgresConfig(cfg *config.Config) *PostgresConfig {
	return &PostgresConfig{
		Host:     cfg.GetEnv("POSTGRES_HOST", "localhost"),
		Port:     cfg.GetEnvAsInt("POSTGRES_PORT", 5432),
		User:     cfg.GetEnv("POSTGRES_USER", "postgres"),
		Password: cfg.GetEnv("POSTGRES_PASSWORD", "postgres"),
		Database: cfg.GetEnv("POSTGRES_DB", "app"),
	}
}


func NewPostgresConnection(lc fx.Lifecycle, cfg *PostgresConfig, logger *zap.Logger) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", 
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	
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
				zap.String("host", cfg.Host), 
				zap.Int("port", cfg.Port),
				zap.String("database", cfg.Database))
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


func NewGoquDatabase(pool *pgxpool.Pool) *goqu.Database {
	return goqu.New("postgres", nil)
}
`

const postgresUserRepositoryTemplate = `package postgres

import (
	"context"
	"errors"
	"time"
	
	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	
	"{{.Name}}/internal/domain"
)


type UserRepository struct {
	pool   *pgxpool.Pool
	db     *goqu.Database
	logger *zap.Logger
}


func NewUserRepository(pool *pgxpool.Pool, db *goqu.Database, logger *zap.Logger) domain.UserRepository {
	return &UserRepository{
		pool:   pool,
		db:     db,
		logger: logger,
	}
}


func (r *UserRepository) Create(user *domain.User) error {
	query, _, err := r.db.Insert("users").
		Rows(goqu.Record{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		}).ToSQL()
	
	if err != nil {
		return err
	}
	
	_, err = r.pool.Exec(context.Background(), query)
	return err
}


func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	query, _, err := r.db.From("users").
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	
	if err != nil {
		return nil, err
	}
	
	var user domain.User
	err = r.pool.QueryRow(context.Background(), query).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}


func (r *UserRepository) List() ([]*domain.User, error) {
	query, _, err := r.db.From("users").ToSQL()
	if err != nil {
		return nil, err
	}
	
	rows, err := r.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	
	return users, nil
}


func (r *UserRepository) Update(user *domain.User) error {
	query, _, err := r.db.Update("users").
		Set(goqu.Record{
			"username":   user.Username,
			"email":      user.Email,
			"updated_at": time.Now(),
		}).
		Where(goqu.C("id").Eq(user.ID)).
		ToSQL()
	
	if err != nil {
		return err
	}
	
	_, err = r.pool.Exec(context.Background(), query)
	return err
}


func (r *UserRepository) Delete(id string) error {
	query, _, err := r.db.Delete("users").
		Where(goqu.C("id").Eq(id)).
		ToSQL()
	
	if err != nil {
		return err
	}
	
	_, err = r.pool.Exec(context.Background(), query)
	return err
}
`

const redisTemplate = `package redis

import (
	"context"
	"fmt"
	
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
	"go.uber.org/zap"
	
	"{{.Name}}/internal/config"
)


type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}


func NewRedisConfig(cfg *config.Config) *RedisConfig {
	return &RedisConfig{
		Host:     cfg.GetEnv("REDIS_HOST", "localhost"),
		Port:     cfg.GetEnvAsInt("REDIS_PORT", 6379),
		Password: cfg.GetEnv("REDIS_PASSWORD", ""),
		DB:       cfg.GetEnvAsInt("REDIS_DB", 0),
	}
}


func NewRedisClient(lc fx.Lifecycle, cfg *RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to Redis", 
				zap.String("host", cfg.Host), 
				zap.Int("port", cfg.Port))
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

const redisUserCacheTemplate = `package redis

import (
	"context"
	"encoding/json"
	"time"
	
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	
	"{{.Name}}/internal/domain"
)


type UserCache struct {
	client *redis.Client
	logger *zap.Logger
	ttl    time.Duration
}


func NewUserCache(client *redis.Client, logger *zap.Logger) *UserCache {
	return &UserCache{
		client: client,
		logger: logger,
		ttl:    time.Hour, 
	}
}


func (c *UserCache) Set(user *domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	
	key := c.userKey(user.ID)
	return c.client.Set(context.Background(), key, data, c.ttl).Err()
}


func (c *UserCache) Get(id string) (*domain.User, error) {
	key := c.userKey(id)
	data, err := c.client.Get(context.Background(), key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil 
		}
		return nil, err
	}
	
	var user domain.User
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, err
	}
	
	return &user, nil
}


func (c *UserCache) Delete(id string) error {
	key := c.userKey(id)
	return c.client.Del(context.Background(), key).Err()
}


func (c *UserCache) userKey(id string) string {
	return "user:" + id
}
`
