package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ruthvik10/membership-managment-system/internal/db/model"
	"github.com/Ruthvik10/membership-managment-system/internal/db/postgres"
	"github.com/Ruthvik10/membership-managment-system/internal/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockStore is a mock implementation of the store interface
type mockStore struct {
	mock.Mock
}

func (m *mockStore) AddMember(ctx context.Context, member *model.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *mockStore) GetMemberByID(ctx context.Context, id uuid.UUID) (*model.Member, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *mockStore) GetMemberByEmail(ctx context.Context, email string) (*model.Member, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Member), args.Error(1)
}

func (m *mockStore) GetAllMembers(ctx context.Context) ([]*model.Member, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Member), args.Error(1)
}

func (m *mockStore) UpdateMember(ctx context.Context, member *model.Member) error {
	args := m.Called(ctx, member)
	return args.Error(0)
}

func (m *mockStore) DeleteMember(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// mockLogger is a mock implementation of the logger interface
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) WriteInfo(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *mockLogger) WriteError(msg string, err error, fields map[string]interface{}) {
	m.Called(msg, err, fields)
}

func (m *mockLogger) WriteFatal(msg string, err error, fields map[string]interface{}) {
	m.Called(msg, err, fields)
}

func TestAddMember(t *testing.T) {
	tests := []struct {
		name           string
		reqBody        addMemberRequest
		setupMock      func(*mocks.MemberStore)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful member creation",
			reqBody: addMemberRequest{
				Name:        "John Doe",
				Email:       "john@example.com",
				PhoneNumber: "1234567890",
				Address:     "123 Main St",
				JoinDate:    time.Now(),
			},
			setupMock: func(ms *mocks.MemberStore) {
				ms.On("AddMember", mock.Anything, mock.MatchedBy(func(m *model.Member) bool {
					return m.Name == "John Doe" && m.Email == "john@example.com"
				})).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "duplicate member",
			reqBody: addMemberRequest{
				Name:        "John Doe",
				Email:       "john@example.com",
				PhoneNumber: "1234567890",
				Address:     "123 Main St",
				JoinDate:    time.Now(),
			},
			setupMock: func(ms *mocks.MemberStore) {
				ms.On("AddMember", mock.Anything, mock.Anything).Return(postgres.ErrMemberAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody: echo.HTTPError{
				Code:    http.StatusConflict,
				Message: "Member already exists",
			},
		},
		{
			name: "missing required fields",
			reqBody: addMemberRequest{
				Name:     "John Doe",
				Address:  "123 Main St",
				JoinDate: time.Now(),
			},
			setupMock:      func(ms *mocks.MemberStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Required fields missing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockStore := new(mocks.MemberStore)
			mockLogger := new(mocks.Logger)
			tt.setupMock(mockStore)

			app := &application{
				store:  mockStore,
				logger: mockLogger,
			}

			// Mock logger calls
			mockLogger.On("WriteError", mock.Anything, mock.Anything, mock.Anything).Maybe()
			mockLogger.On("WriteInfo", mock.Anything, mock.Anything).Maybe()

			// Create request
			reqBody, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/members", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Test
			err := app.addMember(c)

			// Assertions
			if tt.expectedStatus != http.StatusOK && tt.expectedStatus != http.StatusCreated {
				// For error cases, verify the error returned matches expected
				if assert.Error(t, err) {
					httpErr, ok := err.(*echo.HTTPError)
					assert.True(t, ok)
					assert.Equal(t, tt.expectedStatus, httpErr.Code)
					assert.Equal(t, tt.expectedBody.(echo.HTTPError).Message, httpErr.Message)
				}
			} else {
				// For success cases, verify the response
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expectedStatus, rec.Code)
					
					if tt.expectedBody != nil {
						var response interface{}
						err := json.NewDecoder(rec.Body).Decode(&response)
						assert.NoError(t, err)
						assert.Equal(t, tt.expectedBody, response)
					}
				}
			}

			mockStore.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestGetMemberByID(t *testing.T) {
	memberID := uuid.New()
	tests := []struct {
		name           string
		memberID       string
		setupMock      func(*mocks.MemberStore)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:     "successful member retrieval",
			memberID: memberID.String(),
			setupMock: func(ms *mocks.MemberStore) {
				ms.On("GetMemberByID", mock.Anything, memberID).Return(&model.Member{
					ID:          memberID,
					Name:        "John Doe",
					Email:       "john@example.com",
					PhoneNumber: "1234567890",
					Address:     "123 Main St",
					JoinDate:    time.Now(),
					Status:      model.MemberStatusActive,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "member not found",
			memberID: memberID.String(),
			setupMock: func(ms *mocks.MemberStore) {
				ms.On("GetMemberByID", mock.Anything, memberID).Return(nil, postgres.ErrMemberNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Member not found",
			},
		},
		{
			name:           "invalid member ID",
			memberID:       "invalid-uuid",
			setupMock:      func(ms *mocks.MemberStore) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Invalid member ID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockStore := new(mocks.MemberStore)
			mockLogger := new(mocks.Logger)
			tt.setupMock(mockStore)

			app := &application{
				store:  mockStore,
				logger: mockLogger,
			}

			// Mock logger calls
			mockLogger.On("WriteError", mock.Anything, mock.Anything, mock.Anything).Maybe()

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.memberID)

			// Test
			err := app.getMemberByID(c)

			// Assertions
			if tt.expectedStatus != http.StatusOK {
				// For error cases, verify the error returned matches expected
				if assert.Error(t, err) {
					httpErr, ok := err.(*echo.HTTPError)
					assert.True(t, ok)
					assert.Equal(t, tt.expectedStatus, httpErr.Code)
					assert.Equal(t, tt.expectedBody.(echo.HTTPError).Message, httpErr.Message)
				}
			} else {
				// For success cases, verify the response
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expectedStatus, rec.Code)
					
					var response getMemberResponse
					err := json.NewDecoder(rec.Body).Decode(&response)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedBody, response)
				}
			}

			mockStore.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}

func TestGetAllMembers(t *testing.T) {
	// Create test member data
	member1ID := uuid.New()
	member2ID := uuid.New()
	testTime := time.Now()

	tests := []struct {
		name           string
		setupMock      func(*mocks.MemberStore)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful members retrieval",
			setupMock: func(ms *mocks.MemberStore) {
				members := []*model.Member{
					{
						ID:          member1ID,
						Name:        "John Doe",
						Email:       "john@example.com",
						PhoneNumber: "1234567890",
						Address:     "123 Main St",
						JoinDate:    testTime,
						Status:      model.MemberStatusActive,
					},
					{
						ID:          member2ID,
						Name:        "Jane Smith",
						Email:       "jane@example.com",
						PhoneNumber: "0987654321",
						Address:     "456 Oak St",
						JoinDate:    testTime,
						Status:      model.MemberStatusInactive,
					},
				}
				ms.On("GetAllMembers", mock.Anything).Return(members, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: []getMemberResponse{
				{
					ID:          member1ID,
					Name:        "John Doe",
					Email:       "john@example.com",
					PhoneNumber: "1234567890",
					Address:     "123 Main St",
					JoinDate:    testTime,
					Status:      "active",
				},
				{
					ID:          member2ID,
					Name:        "Jane Smith",
					Email:       "jane@example.com",
					PhoneNumber: "0987654321",
					Address:     "456 Oak St",
					JoinDate:    testTime,
					Status:      "inactive",
				},
			},
		},
		{
			name: "database error",
			setupMock: func(ms *mocks.MemberStore) {
				ms.On("GetAllMembers", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get members",
			},
		},
		{
			name: "empty members list",
			setupMock: func(ms *mocks.MemberStore) {
				ms.On("GetAllMembers", mock.Anything).Return([]*model.Member{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []getMemberResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			mockStore := new(mocks.MemberStore)
			mockLogger := new(mocks.Logger)
			tt.setupMock(mockStore)

			// Mock logger calls
			mockLogger.On("WriteError", mock.Anything, mock.Anything, mock.Anything).Maybe()
			mockLogger.On("WriteInfo", mock.Anything, mock.Anything).Maybe()

			app := &application{
				store:  mockStore,
				logger: mockLogger,
			}

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/members", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Test
			err := app.getAllMembers(c)

			// Assertions
			if tt.expectedStatus != http.StatusOK {
				// For error cases, verify the error returned matches expected
				if assert.Error(t, err) {
					httpErr, ok := err.(*echo.HTTPError)
					assert.True(t, ok)
					assert.Equal(t, tt.expectedStatus, httpErr.Code)
					assert.Equal(t, tt.expectedBody.(echo.HTTPError).Message, httpErr.Message)
				}
			} else {
				// For success cases, verify the response
				if assert.NoError(t, err) {
					assert.Equal(t, tt.expectedStatus, rec.Code)
					
					var response []getMemberResponse
					err := json.NewDecoder(rec.Body).Decode(&response)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedBody, response)
				}
			}

			// Verify mocks
			mockStore.AssertExpectations(t)
			mockLogger.AssertExpectations(t)
		})
	}
}
