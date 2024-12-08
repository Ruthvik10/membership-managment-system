package mocks

import (
	"context"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MemberStore is a mock implementation of the store interface
type MemberStore struct {
	mock.Mock
}

func (m *MemberStore) AddMember(ctx context.Context, member *model.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MemberStore) GetMemberByID(ctx context.Context, id uuid.UUID) (*model.Member, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *MemberStore) GetMemberByEmail(ctx context.Context, email string) (*model.Member, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *MemberStore) GetAllMembers(ctx context.Context) ([]*model.Member, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Member), args.Error(1)
}

func (m *MemberStore) UpdateMember(ctx context.Context, member *model.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *MemberStore) DeleteMember(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
