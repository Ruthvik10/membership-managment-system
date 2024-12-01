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

type MemberStore struct {
	conn *pgxpool.Pool
}

func NewMemberStore(conn *pgxpool.Pool) *MemberStore {
	return &MemberStore{
		conn: conn,
	}
}

func (s *MemberStore) AddMember(ctx context.Context, member *model.Member) error {
	query := `
		INSERT INTO members (name, email, phone, address, join_date, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`
	args := []any{
		member.Name,
		member.Email,
		member.PhoneNumber,
		member.Address,
		member.JoinDate,
		member.Status,
	}

	if err := s.conn.QueryRow(ctx, query, args...).Scan(
		&member.ID,
		&member.Name,
		&member.Email,
		&member.PhoneNumber,
		&member.Address,
		&member.JoinDate,
		&member.Status,
	); err != nil {
		switch {
		case IsPgError(err, PgUniqueViolation):
			return fmt.Errorf("%w: %v", ErrMemberAlreadyExists, err)
		case IsPgError(err, PgNotNullViolation):
			return fmt.Errorf("%w: %v", ErrMissingRequiredField, err)
		default:
			return fmt.Errorf("failed to add member: %w", err)
		}
	}
	return nil
}

func (s *MemberStore) GetMemberByID(ctx context.Context, id uuid.UUID) (*model.Member, error) {
	query := `
		SELECT id, name, email, phone, address, join_date, status
		FROM members
		WHERE id = $1
	`
	var member model.Member
	if err := s.conn.QueryRow(ctx, query, id).Scan(
		&member.ID,
		&member.Name,
		&member.Email,
		&member.PhoneNumber,
		&member.Address,
		&member.JoinDate,
		&member.Status,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("%w: %w", ErrMemberNotFound, err)
		default:
			return nil, fmt.Errorf("failed to get member: %w", err)
		}
	}
	return &member, nil
}

func (s *MemberStore) GetMemberByEmail(ctx context.Context, email string) (*model.Member, error) {
	query := `
		SELECT id, name, email, phone, address, join_date, status
		FROM members
		WHERE email = $1
	`
	var member model.Member
	if err := s.conn.QueryRow(ctx, query, email).Scan(
		&member.ID,
		&member.Name,
		&member.Email,
		&member.PhoneNumber,
		&member.Address,
		&member.JoinDate,
		&member.Status,
	); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("%w: %w", ErrMemberNotFound, err)
		default:
			return nil, fmt.Errorf("failed to get member: %w", err)
		}
	}
	return &member, nil
}

func (s *MemberStore) GetAllMembers(ctx context.Context) ([]*model.Member, error) {
	query := `
		SELECT id, name, email, phone, address, join_date, status
		FROM members
	`
	rows, err := s.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}
	defer rows.Close()

	var members []*model.Member
	for rows.Next() {
		var member model.Member
		if err := rows.Scan(
			&member.ID,
			&member.Name,
			&member.Email,
			&member.PhoneNumber,
			&member.Address,
			&member.JoinDate,
			&member.Status,
		); err != nil {
			return nil, fmt.Errorf("failed to scan member: %w", err)
		}
		members = append(members, &member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over members: %w", err)
	}
	return members, nil
}

func (s *MemberStore) UpdateMember(ctx context.Context, member *model.Member) error {
	query := `
		UPDATE members
		SET name = $1, email = $2, phone = $3, address = $4, join_date = $5, status = $6
		WHERE id = $7
		RETURNING *
	`
	args := []any{
		member.Name,
		member.Email,
		member.PhoneNumber,
		member.Address,
		member.JoinDate,
		member.Status,
		member.ID,
	}
	err := s.conn.QueryRow(ctx, query, args...).Scan(
		&member.ID,
		&member.Name,
		&member.Email,
		&member.PhoneNumber,
		&member.Address,
		&member.JoinDate,
		&member.Status,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return fmt.Errorf("%w: %w", ErrMemberNotFound, err)
		case IsPgError(err, PgNotNullViolation):
			return fmt.Errorf("%w: %v", ErrMissingRequiredField, err)
		case IsPgError(err, PgUniqueViolation):
			return fmt.Errorf("%w: %v", ErrMemberAlreadyExists, err)
		default:
			return fmt.Errorf("failed to update member: %w", err)
		}
	}
	return nil
}

func (s *MemberStore) DeleteMember(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM members
		WHERE id = $1
	`
	rows, err := s.conn.Exec(ctx, query, id)
	if rows.RowsAffected() == 0 {
		return fmt.Errorf("%w: %w", ErrMemberNotFound, err)
	}
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}
