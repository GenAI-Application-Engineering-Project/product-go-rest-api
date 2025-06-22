package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSystemUtil struct {
	mock.Mock
}

func (m *MockSystemUtil) CurrentTime() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

func (m *MockSystemUtil) NewUUID() uuid.UUID {
	args := m.Called()
	return args.Get(0).(uuid.UUID)
}
