package storer

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

var _ Storer = &mockStorer{}

type mockStorer struct {
	mock.Mock
}

func NewMock() *mockStorer {
	return &mockStorer{}
}

func (s *mockStorer) RetrieveIndividualFiles(reportType ReportType, year int, month int) ([][]byte, error) {
	args := s.Called(reportType, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([][]byte), args.Error(1)
}

func (s *mockStorer) StubRetrieveIndividualFiles(reportType interface{}, year interface{}, month interface{}) *mock.Call {
	return s.On("RetrieveIndividualFiles", reportType, year, month)
}

func (s *mockStorer) AssertRetrieveIndividualFilesNotCalled(t *testing.T) {
	s.AssertNotCalled(t, "RetrieveIndividualFiles", mock.Anything, mock.Anything, mock.Anything)
}

func (s *mockStorer) RetrieveAggregated(reportType ReportType, year int, month int) ([]byte, error) {
	args := s.Called(reportType, year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]byte), args.Error(1)
}

func (s *mockStorer) StubRetrieveAggregated(reportType interface{}, year interface{}, month interface{}) *mock.Call {
	return s.On("RetrieveAggregated", reportType, year, month)
}

func (s *mockStorer) StoreAggregated(reportType ReportType, year int, month int, data []byte) error {
	args := s.Called(reportType, year, month, data)
	return args.Error(0)
}

func (s *mockStorer) StubStoreAggregated(reportType interface{}, year interface{}, month interface{}, data interface{}) *mock.Call {
	return s.On("StoreAggregated", reportType, year, month, data)
}

func (s *mockStorer) AssertStoreAggregatedCalled(t *testing.T, reportType interface{}, year interface{}, month interface{}, data interface{}) {
	s.AssertCalled(t, "StoreAggregated", reportType, year, month, data)
}

func (s *mockStorer) AssertStoreAggregatedNotCalled(t *testing.T) {
	s.AssertNotCalled(t, "StoreAggregated", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}
