package mocks

import "github.com/stretchr/testify/mock"

// Logger is a mock implementation of the logger interface
type Logger struct {
	mock.Mock
}

func (m *Logger) WriteInfo(msg string, fields map[string]interface{}) {
	m.Called(msg, fields)
}

func (m *Logger) WriteError(msg string, err error, fields map[string]interface{}) {
	m.Called(msg, err, fields)
}

func (m *Logger) WriteFatal(msg string, err error, fields map[string]interface{}) {
	m.Called(msg, err, fields)
}
