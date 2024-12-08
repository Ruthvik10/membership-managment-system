package model

import "github.com/google/uuid"

type Sport struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
}

func (s *Sport) Valid() bool {
	if len(s.Name) <= 2 {
		return false
	}
	return true
}
