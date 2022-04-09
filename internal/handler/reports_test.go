//go:build unit

package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hpcsc/outside-in-go/internal/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReports(t *testing.T) {
	t.Run("return 404 when report type route parameter is neither single nor cumulative", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/reports/not-valid", nil)
		require.NoError(t, err)
		recorder := httptest.NewRecorder()
		stubGenerator := report.NewMockGenerator()
		router := testRouterWithReports(stubGenerator)

		router.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNotFound, recorder.Code)
	})

	t.Run("single", func(t *testing.T) {
		t.Run("return 200 with csv file when report is generated successfully", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/reports/single?year=2022&month=4", nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			stubGenerator := report.NewMockGenerator()
			stubGenerator.StubGenerateSingle().Return([]byte("some,csv,data"), nil)
			router := testRouterWithReports(stubGenerator)

			router.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, "text/csv", recorder.Header().Get("Content-Type"))
			assert.Equal(t, "attachment; filename=single-202204.csv", recorder.Header().Get("Content-Disposition"))
			assert.Equal(t, "some,csv,data", recorder.Body.String())
		})
	})
}

func testRouterWithReports(generator report.Generator) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	RegisterReportsRoutes(r, generator)
	return r
}
