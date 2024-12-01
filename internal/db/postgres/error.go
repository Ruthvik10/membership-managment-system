package postgres

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrMemberAlreadyExists   = errors.New("member already exists")
	ErrMemberNotFound        = errors.New("member not found")
	ErrMissingRequiredField  = errors.New("missing required field")
)

const (
	// PostgreSQL error codes
	PgUniqueViolation    = "23505" // unique_violation
	PgNotNullViolation   = "23502" // not_null_violation
	PgForeignKeyViolation = "23503" // foreign_key_violation
)

// IsPgError checks if the error is a PostgreSQL error with the given code
func IsPgError(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}
