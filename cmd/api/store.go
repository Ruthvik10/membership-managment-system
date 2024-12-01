package main

import (
	"context"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/google/uuid"
)

type memberStore interface {
	AddMember(ctx context.Context, member *model.Member) error
	GetMemberByID(ctx context.Context, id uuid.UUID) (*model.Member, error)
	GetMemberByEmail(ctx context.Context, email string) (*model.Member, error)
	GetAllMembers(ctx context.Context) ([]*model.Member, error)
	UpdateMember(ctx context.Context, member *model.Member) error
	DeleteMember(ctx context.Context, id uuid.UUID) error
}

type store interface {
	memberStore
}
