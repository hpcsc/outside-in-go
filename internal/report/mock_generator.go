package report

import (
	"github.com/stretchr/testify/mock"
)

var _ Generator = &mockGenerator{}

type mockGenerator struct {
	mock.Mock
}

func NewMockGenerator() *mockGenerator {
	return &mockGenerator{}
}

func (m *mockGenerator) GenerateSingle(year int, month int) ([]byte, error) {
	args := m.Called(year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockGenerator) GenerateCumulative(year int, month int) ([]byte, error) {
	args := m.Called(year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}
