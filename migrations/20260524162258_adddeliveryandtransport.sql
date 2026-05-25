-- +goose Up
-- 1. Добавляем новую колонку в существующую таблицу
ALTER TABLE couriers ADD COLUMN transport_type TEXT NOT NULL DEFAULT 'on_foot';

-- 2. Создаем новую таблицу
CREATE TABLE delivery (
    id                  BIGSERIAL PRIMARY KEY,
    courier_id          BIGINT NOT NULL,
    order_id            VARCHAR(255) NOT NULL,
    assigned_at         TIMESTAMP NOT NULL DEFAULT NOW(),
    deadline            TIMESTAMP NOT NULL
);

-- +goose Down
-- Откатываем всё в обратном порядке:
DROP TABLE delivery;
ALTER TABLE couriers DROP COLUMN transport_type;