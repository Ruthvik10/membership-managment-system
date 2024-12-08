package mocks

import (
	"context"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SportStore struct {
	mock.Mock
}

func (s *MemberStore) GetAllSports(ctx context.Context) ([]*model.Sport, error) {
	args := s.Called(ctx)
	return args.Get(0).([]*model.Sport), args.Error(1)
}

func (s *MemberStore) GetSportByID(ctx context.Context, id uuid.UUID) (*model.Sport, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*model.Sport), args.Error(1)
}

func (s *MemberStore) AddSport(ctx context.Context, sport *model.Sport) error {
	args := s.Called(ctx, sport)
	return args.Error(0)
}

func (s *MemberStore) UpdateSport(ctx context.Context, sport *model.Sport) error {
	args := s.Called(ctx, sport)
	return args.Error(0)
}

func (s *MemberStore) DeleteSport(ctx context.Context, id uuid.UUID) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}
