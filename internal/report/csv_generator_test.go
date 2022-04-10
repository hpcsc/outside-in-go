//go:build unit

package report

import (
	"errors"
	"github.com/hpcsc/outside-in-go/internal/storer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCsvGenerator(t *testing.T) {
	t.Run("generate single", func(t *testing.T) {
		generateReportTestSuite(t, storer.SingleReportType, func(gen Generator, year int, month int) ([]byte, error) {
			return gen.GenerateSingle(year, month)
		})
	})

	t.Run("generate cumulative", func(t *testing.T) {
		generateReportTestSuite(t, storer.CumulativeReportType, func(gen Generator, year int, month int) ([]byte, error) {
			return gen.GenerateCumulative(year, month)
		})
	})
}

func generateReportTestSuite(
	t *testing.T,
	reportType storer.ReportType,
	generate func(gen Generator, year int, month int) ([]byte, error),
) {
	year := 2022
	month := 4

	t.Run("return aggregated file straightaway when found", func(t *testing.T) {
		stubStorer := storer.NewMock()
		stubStorer.StubRetrieveAggregated(reportType, year, month).Return([]byte("some,data"), nil)
		g := NewCsvGenerator(stubStorer)

		data, err := generate(g, year, month)

		require.NoError(t, err)
		require.Equal(t, []byte("some,data"), data)
		stubStorer.AssertRetrieveIndividualFilesNotCalled(t)
		stubStorer.AssertStoreAggregatedNotCalled(t)
	})

	t.Run("return error if no individual files available for given year and month", func(t *testing.T) {
		stubStorer := storer.NewMock()
		stubStorer.StubRetrieveAggregated(reportType, year, month).Return(nil, nil)
		stubStorer.StubRetrieveIndividualFiles(reportType, year, month).Return([][]byte{}, nil)
		g := NewCsvGenerator(stubStorer)

		_, err := generate(g, year, month)

		require.Error(t, err)
		require.Contains(t, err.Error(), "no data available for 04/2022")
	})

	t.Run("return error if failed to retrieve individual files", func(t *testing.T) {
		stubStorer := storer.NewMock()
		stubStorer.StubRetrieveAggregated(reportType, year, month).Return(nil, nil)
		stubStorer.StubRetrieveIndividualFiles(reportType, year, month).Return(nil, errors.New("some error"))
		g := NewCsvGenerator(stubStorer)

		_, err := generate(g, year, month)

		require.Error(t, err)
		require.Contains(t, err.Error(), "some error")
	})

	t.Run("store and return aggregated report", func(t *testing.T) {
		stubStorer := storer.NewMock()
		stubStorer.StubRetrieveAggregated(reportType, year, month).Return(nil, nil)
		stubStorer.StubRetrieveIndividualFiles(reportType, year, month).Return(
			[][]byte{
				[]byte(`CLUSTER,DATA
cluster-1,data-1.1
cluster-1,data-1.2`),
				[]byte(`CLUSTER,DATA
cluster-2,data-2.1
cluster-2,data-2.2`),
			},
			nil)
		stubStorer.StubStoreAggregated(reportType, year, month, mock.Anything).Return(nil)
		g := NewCsvGenerator(stubStorer)

		data, err := generate(g, year, month)

		require.NoError(t, err)
		stubStorer.AssertStoreAggregatedCalled(t, reportType, year, month, mock.Anything)
		expecteData := []byte(`CLUSTER,DATA
cluster-1,data-1.1
cluster-1,data-1.2
cluster-2,data-2.1
cluster-2,data-2.2
`)
		require.Equal(t, expecteData, data)
	})

	t.Run("return error when failed to store aggregate file", func(t *testing.T) {
		stubStorer := storer.NewMock()
		stubStorer.StubRetrieveAggregated(reportType, year, month).Return(nil, nil)
		stubStorer.StubRetrieveIndividualFiles(reportType, year, month).Return(
			[][]byte{
				[]byte(`CLUSTER,DATA
cluster-1,data-1.1
cluster-1,data-1.2`),
			},
			nil)
		stubStorer.StubStoreAggregated(reportType, year, month, mock.Anything).Return(errors.New("some error"))
		g := NewCsvGenerator(stubStorer)

		_, err := generate(g, year, month)

		require.Error(t, err)
		require.Contains(t, err.Error(), "some error")
	})
}
