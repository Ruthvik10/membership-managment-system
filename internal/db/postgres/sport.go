package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SportStore struct {
	conn *pgxpool.Pool
}

func NewSportStore(conn *pgxpool.Pool) *SportStore {
	return &SportStore{
		conn: conn,
	}
}

func (s *SportStore) AddSport(ctx context.Context, sport *model.Sport) error {
	query := `
		INSERT INTO sports (name, description)
		VALUES ($1, $2) RETURNING *
	`
	args := []any{sport.Name, sport.Description}
	if err := s.conn.QueryRow(ctx, query, args...).Scan(
		&sport.ID,
		&sport.Name,
		&sport.Description,
	); err != nil {
		switch {
		case IsPgError(err, PgUniqueViolation):
			return fmt.Errorf("%w: %w", ErrSportAlreadyExists, err)
		case IsPgError(err, PgNotNullViolation):
			return fmt.Errorf("%w: %w", ErrMissingRequiredField, err)
		default:
			return fmt.Errorf("failed to add sport: %w", err)
		}
	}
	return nil
}

func (s *SportStore) GetSportByID(ctx context.Context, id uuid.UUID) (*model.Sport, error) {
	query := `
		SELECT id, name, description
		FROM sports
		WHERE id = $1
	`
	var sport model.Sport
	if err := s.conn.QueryRow(ctx, query, id).Scan(
		&sport.ID,
		&sport.Name,
		&sport.Description,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("%w: %w", ErrSportNotFound, err)
		default:
			return nil, fmt.Errorf("failed to get sport: %w", err)
		}
	}
	return &sport, nil
}

func (s *SportStore) GetAllSports(ctx context.Context) ([]*model.Sport, error) {
	query := `
		SELECT id, name, description
		FROM sports
	`
	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get sports: %w", err)
	}
	defer rows.Close()

	var sports []*model.Sport
	for rows.Next() {
		var sport model.Sport
		if err := rows.Scan(
			&sport.ID,
			&sport.Name,
			&sport.Description,
		); err != nil {
			return nil, fmt.Errorf("failed to scan sport: %w", err)
		}
		sports = append(sports, &sport)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over sports: %w", err)
	}
	return sports, nil
}

func (s *SportStore) DeleteSport(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM sports
		WHERE id = $1
	`
	rows, err := s.conn.Exec(ctx, query, id)
	if rows.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", ErrSportNotFound, err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete sport: %w", err)
	}
	return nil
}

func (s *SportStore) UpdateSport(ctx context.Context, sport *model.Sport) error {
	query := `
		UPDATE sports
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING *
	`
	args := []any{
		sport.Name,
		sport.Description,
		sport.ID,
	}
	if err := s.conn.QueryRow(ctx, query, args...).Scan(
		&sport.ID,
		&sport.Name,
		&sport.Description,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return fmt.Errorf("%w: %w", ErrSportNotFound, err)
		case IsPgError(err, PgNotNullViolation):
			return fmt.Errorf("%w: %w", ErrMissingRequiredField, err)
		case IsPgError(err, PgUniqueViolation):
			return fmt.Errorf("%w: %w", ErrSportAlreadyExists, err)
		default:
			return fmt.Errorf("failed to update sport: %w", err)
		}
	}
	return nil
}
