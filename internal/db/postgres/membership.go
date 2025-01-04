package postgres

import (
	"context"
	"fmt"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MembershipStore struct {
	conn *pgxpool.Pool
}

func NewMembershipStore(conn *pgxpool.Pool) *MembershipStore {
	return &MembershipStore{
		conn: conn,
	}
}

func (s *MembershipStore) AddMembership(ctx context.Context, membership *model.Membership) error {
	query := `INSERT INTO memberships (member_id, sport_id, type, start_date, due_date, status, fee) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *`
	args := []any{
		membership.MemberID,
		membership.SportID,
		membership.Type,
		membership.StartDate,
		membership.DueDate,
		membership.Status,
		membership.Fee,
	}
	err := s.conn.QueryRow(ctx, query, args...).Scan(
		&membership.ID,
		&membership.MemberID,
		&membership.SportID,
		&membership.Type,
		&membership.StartDate,
		&membership.DueDate,
		&membership.Status,
		&membership.Fee,
	)
	if err != nil {
		switch {
		case IsPgError(err, PgUniqueViolation):
			return fmt.Errorf("%w: %w", ErrMembershipAlreadyExists, err)
		case IsPgError(err, PgNotNullViolation):
			return fmt.Errorf("%w: %w", ErrMissingRequiredField, err)
		default:
			return fmt.Errorf("failed to add membership: %w", err)
		}
	}
	return nil
}
