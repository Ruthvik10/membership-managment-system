package model

import (
	"regexp"

	"github.com/google/uuid"
)

type MemberStatus int

const (
	MemberStatusInactive MemberStatus = iota
	MemberStatusActive
)

var MemberStatusMap = map[MemberStatus]string{
	MemberStatusInactive: "Inactive",
	MemberStatusActive:   "Active",
}

type Member struct {
	ID          uuid.UUID    `db:"id"`
	Name        string       `db:"name"`
	Email       string       `db:"email"`
	PhoneNumber string       `db:"phone"`
	Address     string       `db:"address"`
	Status      MemberStatus `db:"status"`
}

// TODO: Return which fields are invalid
func (m *Member) Valid() bool {
	if len(m.Name) <= 2 {
		return false
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(m.Email) {
		return false
	}

	if len(m.PhoneNumber) != 10 {
		return false
	}

	if m.Status != MemberStatusInactive && m.Status != MemberStatusActive {
		return false
	}

	return true
}
