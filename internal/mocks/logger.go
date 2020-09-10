package mocks

import (
	"context"

	"github.com/redselig/currencier/internal/domain/usecase"
)

var _ usecase.Logger = (*MockLogger)(nil)

type MockLogger struct{}

func NewMockLogger() *MockLogger {
	return &MockLogger{}
}
func (m MockLogger) Log(ctx context.Context, message interface{}, args ...interface{}) {
}
