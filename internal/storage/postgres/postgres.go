package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, connectionString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}
	config.MaxConns = 20                            // Максимальное количество соединений
	config.MinConns = 5                             // Минимальное количество соединений
	config.MaxConnLifetime = time.Hour              // Время жизни соединения
	config.MaxConnIdleTime = time.Minute * 30       // Время бездействия
	config.HealthCheckPeriod = time.Minute          // Период проверки соединения
	pool, err := pgxpool.NewWithConfig(ctx, config) // Создаем пул
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx) // Проверяем подключение
	if err != nil {
		return nil, err
	}
	return pool, nil
}
