package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/doter2410/avito-project/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CourierPostgres struct {
	pool *pgxpool.Pool
}

func NewCourierPostgres(pool *pgxpool.Pool) *CourierPostgres {
	return &CourierPostgres{pool: pool}
}

func (s *CourierPostgres) CreateCourier(ctx context.Context, c model.Courier) (int64, error) {
	query := `INSERT INTO couriers (name, phone, status, transport_type) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int64
	err := s.pool.QueryRow(ctx, query, c.Name, c.Phone, c.Status, c.TransportType).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *CourierPostgres) GetCourierById(ctx context.Context, id int64) (*model.Courier, error) {
	query := `SELECT id, name, phone, status, created_at, updated_at, transport_type FROM couriers WHERE id = $1`
	var c model.Courier
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Phone,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.TransportType,
	)

	if err != nil {
		return nil, err
	}
	return &c, err
}

func (s *CourierPostgres) GetCouriers(ctx context.Context) ([]*model.Courier, error) {
	query := `SELECT id, name, phone, status, created_at, updated_at, transport_type FROM couriers`
	var res []*model.Courier
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := model.Courier{}
		err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt, &c.TransportType)
		if err != nil {
			return nil, err
		}
		res = append(res, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if res == nil {
		res = []*model.Courier{}
	}

	return res, nil

}

func (s *CourierPostgres) UpdateCourier(ctx context.Context, id int64, c model.Courier) error {
	query := `UPDATE couriers SET name = $1, phone = $2, status = $3, updated_at = now() WHERE id = $4`

	res, err := s.pool.Exec(ctx, query, c.Name, c.Phone, c.Status, id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("courier with id &d not found", id)
	}
	return nil
}

func (r *CourierPostgres) AssignCourierToOrder(ctx context.Context, orderID string, calcDeadline func(transportType string) time.Time) (*model.Courier, *model.Delivery, error) {
	// 1. Открываем транзакцию
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	// Если мы забудем сделать Commit, функция завершится и Rollback всё отменит. Это наша страховка.
	defer tx.Rollback(ctx)

	// 2. Ищем свободного курьера и БЛОКИРУЕМ эту строку (чтобы другой заказ его не перехватил)
	var c model.Courier
	// ТВОЯ ЗАДАЧА: Написать SELECT запрос, который вернет id, name, phone, status, transport_type
	// из таблицы couriers, где status = 'available', взять только 1 запись (LIMIT 1)
	// ВАЖНО: В конце SQL-запроса добавь слова FOR UPDATE (это заблокирует строку до конца транзакции)
	querySelect := `SELECT id, name, phone, status, transport_type FROM couriers WHERE status = 'available' LIMIT 1 FOR UPDATE`
	err = tx.QueryRow(ctx, querySelect).Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.TransportType)
	if err != nil {
		return nil, nil, err // Сюда мы попадем, если свободных курьеров нет
	}

	// 3. Высчитываем дедлайн с помощью переданной функции из Фабрики!
	deadline := calcDeadline(c.TransportType)

	// 4. Обновляем статус курьера
	// ТВОЯ ЗАДАЧА: Написать UPDATE запрос, который ставит status = 'busy' для курьера с id = c.ID
	queryUpd := `UPDATE couriers SET status = 'busy' WHERE id = $1`

	_, err = tx.Exec(ctx, queryUpd, c.ID)
	if err != nil {
		return nil, nil, err
	}

	// 5. Записываем заказ в таблицу delivery
	var d model.Delivery
	d.CourierID = c.ID
	d.OrderID = orderID
	d.Deadline = deadline
	// ТВОЯ ЗАДАЧА: Написать INSERT запрос в таблицу delivery.
	// Верни сгенерированный id и assigned_at (используй RETURNING id, assigned_at)
	queryInsert := `INSERT INTO delivery (courier_id, order_id, deadline) VALUES ($1, $2, $3) RETURNING id, assigned_at`
	err = tx.QueryRow(ctx, queryInsert, d.CourierID, d.OrderID, d.Deadline).Scan(&d.ID, &d.AssignedAt)
	if err != nil {
		return nil, nil, err
	}

	// 6. Подтверждаем транзакцию! Без этого в базе ничего не сохранится.
	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	// Обновляем структуру для возврата
	c.Status = "busy"
	return &c, &d, nil
}

// 1. Меняем возвращаемое значение на (*model.Courier, error), как в интерфейсе
func (r *CourierPostgres) UnassignCourierFromOrder(ctx context.Context, orderID string) (*model.Courier, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err // Возвращаем nil вместо структуры
	}
	defer tx.Rollback(ctx)

	var courierID int64
	querySearch := `SELECT courier_id FROM delivery WHERE order_id = $1 FOR UPDATE`
	// Тут передаем orderID вторым аргументом (ты забыл его передать в своем коде)
	err = tx.QueryRow(ctx, querySearch, orderID).Scan(&courierID)
	if err != nil {
		return nil, err
	}

	queryDelete := `DELETE FROM delivery WHERE order_id = $1`
	_, err = tx.Exec(ctx, queryDelete, orderID)
	if err != nil {
		return nil, err
	}

	queryUpd := `UPDATE couriers SET status = 'available' WHERE id = $1`
	_, err = tx.Exec(ctx, queryUpd, courierID)
	if err != nil { // Добавили проверку ошибки!
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// 2. Возвращаем структуру курьера с заполненным ID
	return &model.Courier{ID: courierID}, nil
}
