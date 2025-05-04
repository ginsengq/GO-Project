package carrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"myproject/internal/entities"

	"github.com/jackc/pgx/v4/pgxpool"
)

type postgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) Repository {
	return &postgresRepo{db: db}
}

func (r *postgresRepo) Create(ctx context.Context, car *entities.Car) (int, error) {
	query := `
		INSERT INTO cars (brand, model, year, price, mileage, color, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id`

	var id int
	err := r.db.QueryRow(ctx, query,
		car.Brand,
		car.Model,
		car.Year,
		car.Price,
		car.Mileage,
		car.Color,
		car.Status,
	).Scan(&id)
	return id, err
}

func (r *postgresRepo) GetByID(ctx context.Context, id int) (*entities.Car, error) {
	var car entities.Car
	query := `
		SELECT id, brand, model, year, price, mileage, color, status, created_at, updated_at
		FROM cars WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&car.ID, &car.Brand, &car.Model, &car.Year,
		&car.Price, &car.Mileage, &car.Color, &car.Status,
		&car.CreatedAt, &car.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &car, err
}

func (r *postgresRepo) Update(ctx context.Context, id int, update entities.CarUpdate) error {
	var sets []string
	var args []interface{}
	argPos := 1

	if update.Brand != nil {
		sets = append(sets, fmt.Sprintf("brand = $%d", argPos))
		args = append(args, *update.Brand)
		argPos++
	}
	if update.Model != nil {
		sets = append(sets, fmt.Sprintf("model = $%d", argPos))
		args = append(args, *update.Model)
		argPos++
	}
	if update.Year != nil {
		sets = append(sets, fmt.Sprintf("year = $%d", argPos))
		args = append(args, *update.Year)
		argPos++
	}
	if update.Price != nil {
		sets = append(sets, fmt.Sprintf("price = $%d", argPos))
		args = append(args, *update.Price)
		argPos++
	}
	if update.Mileage != nil {
		sets = append(sets, fmt.Sprintf("mileage = $%d", argPos))
		args = append(args, *update.Mileage)
		argPos++
	}
	if update.Color != nil {
		sets = append(sets, fmt.Sprintf("color = $%d", argPos))
		args = append(args, *update.Color)
		argPos++
	}
	if update.Status != nil {
		sets = append(sets, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *update.Status)
		argPos++
	}

	if len(sets) == 0 {
		return nil // No fields to update
	}

	query := fmt.Sprintf("UPDATE cars SET %s, updated_at = NOW() WHERE id = $%d", strings.Join(sets, ", "), argPos)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *postgresRepo) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM cars
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresRepo) List(ctx context.Context, filter entities.CarFilter) ([]*entities.Car, error) {
	var cars []*entities.Car
	var whereClauses []string
	var args []interface{}
	argPos := 1

	if filter.Brand != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("brand = $%d", argPos))
		args = append(args, *filter.Brand)
		argPos++
	}
	if filter.Model != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("model = $%d", argPos))
		args = append(args, *filter.Model)
		argPos++
	}
	if filter.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price >= $%d", argPos))
		args = append(args, *filter.MinPrice)
		argPos++
	}
	if filter.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argPos))
		args = append(args, *filter.MaxPrice)
		argPos++
	}
	if filter.YearFrom != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("year >= $%d", argPos))
		args = append(args, *filter.YearFrom)
		argPos++
	}
	if filter.YearTo != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("year <= $%d", argPos))
		args = append(args, *filter.YearTo)
		argPos++
	}
	if filter.Status != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *filter.Status)
		argPos++
	}
	if filter.Color != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("color = $%d", argPos))
		args = append(args, *filter.Color)
		argPos++
	}

	query := `
		SELECT id, brand, model, year, price, mileage, color, status, created_at, updated_at
		FROM cars`

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	if filter.SortBy != nil {
		sortOrder := "ASC"
		if filter.SortOrder != nil && strings.ToUpper(*filter.SortOrder) == "DESC" {
			sortOrder = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", *filter.SortBy, sortOrder)
	}

	if filter.Limit != nil {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
	}
	if filter.Offset != nil {
		query += fmt.Sprintf(" OFFSET %d", filter.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car entities.Car
		if err := rows.Scan(
			&car.ID, &car.Brand, &car.Model, &car.Year,
			&car.Price, &car.Mileage, &car.Color, &car.Status,
			&car.CreatedAt, &car.UpdatedAt,
		); err != nil {
			return nil, err
		}
		cars = append(cars, &car)
	}

	return cars, nil
}

func (r *postgresRepo) SetStatus(ctx context.Context, id int, status string) error {
	query := `
		UPDATE cars
		SET status = $2, updated_at = NOW()
		WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id, status)
	return err
}

var (
	ErrNotFound      = errors.New("car not found")
	ErrInvalidID     = errors.New("invalid car ID")
	ErrInvalidStatus = errors.New("invalid car status")
)
