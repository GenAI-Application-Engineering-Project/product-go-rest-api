package mocks

import (
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
)

type MockAppLogger struct {
	mock.Mock
}

func (m *MockAppLogger) Logger() zerolog.Logger {
	args := m.Called()
	return args.Get(0).(zerolog.Logger)
}

func (m *MockAppLogger) Fatal(err error, msg string) {
	m.Called(err, msg)
}
