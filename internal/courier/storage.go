package courier

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(pool *pgxpool.Pool) *Storage {
	return &Storage{pool: pool}
}

func (s *Storage) CreateCourier(ctx context.Context, c Courier) (int64, error) {
	query := `INSERT INTO couriers (name, phone, status) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := s.pool.QueryRow(ctx, query, c.Name, c.Phone, c.Status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *Storage) GetCourierById(ctx context.Context, id int64) (*Courier, error) {
	query := `SELECT id, name, phone, status, created_at, updated_at FROM couriers WHERE id = $1`
	var c Courier
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&c.ID,
		&c.Name,
		&c.Phone,
		&c.Status,
		&c.CreatedAt,
		&c.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &c, err
}

func (s *Storage) GetCouriers(ctx context.Context) ([]*Courier, error) {
	query := `SELECT id, name, phone, status, created_at, updated_at FROM couriers`
	var res []*Courier
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		c := Courier{}
		err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if res == nil {
		res = []*Courier{}
	}

	return res, nil

}

func (s *Storage) UpdateCourier(ctx context.Context, id int64, c Courier) error {
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
