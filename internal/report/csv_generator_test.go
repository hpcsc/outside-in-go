//go:build unit

package report

import (
	"github.com/hpcsc/outside-in-go/internal/storer"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCsvGenerator(t *testing.T) {
	t.Run("generate single", func(t *testing.T) {
		t.Run("return aggregated file straightaway when found", func(t *testing.T) {
			year := 2022
			month := 4
			stubStorer := storer.NewMock()
			stubStorer.StubRetrieveAggregated(storer.SingleReportType, year, month).Return([]byte("some,data"), nil)
			g := NewCsvGenerator(stubStorer)

			data, err := g.GenerateSingle(year, month)

			require.NoError(t, err)
			require.Equal(t, []byte("some,data"), data)
			stubStorer.AssertRetrieveIndividualFilesNotCalled(t)
			stubStorer.AssertStoreAggregatedNotCalled(t)
		})
	})
}
