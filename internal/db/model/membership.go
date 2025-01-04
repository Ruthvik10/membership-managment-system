package model

import (
	"time"

	"github.com/google/uuid"
)

type MembershipStatus int

var (
	MembershipActive   MembershipStatus = 1
	MembershipInactive MembershipStatus = 0
)

type MembershipType string

var (
	MembershipTypeMembership MembershipType = "membership"
	MembershipTypeTraining   MembershipType = "training"
)

type Membership struct {
	ID        uuid.UUID        `db:"id"`
	MemberID  uuid.UUID        `db:"member_id"`
	SportID   uuid.UUID        `db:"sport_id"`
	Type      MembershipType   `db:"type"`
	StartDate time.Time        `db:"start_date"`
	DueDate   time.Time        `db:"due_date"`
	Status    MembershipStatus `db:"status"`
	Fee       float64          `db:"fee"`
}

func (m *Membership) Valid() bool {
	if m.Type != MembershipTypeMembership && m.Type != MembershipTypeTraining {
		return false
	}

	if m.Fee <= 0 {
		return false
	}

	if m.Status != MembershipActive && m.Status != MembershipInactive {
		return false
	}

	if m.StartDate.After(m.DueDate) {
		return false
	}

	return true
}
