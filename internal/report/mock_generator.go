package report

var _ Generator = &mockGenerator{}

type mockGenerator struct {
}

func NewMockGenerator() *mockGenerator {
	return &mockGenerator{}
}
