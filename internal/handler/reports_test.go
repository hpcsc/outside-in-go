//go:build unit

package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hpcsc/outside-in-go/internal/report"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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

		t.Run("set year and month to previous month if both parameters are absent", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/reports/single", nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			stubGenerator := report.NewMockGenerator()
			stubGenerator.StubGenerateSingle().Return([]byte("some,csv,data"), nil)
			router := testRouterWithReports(stubGenerator)

			router.ServeHTTP(recorder, req)

			now := time.Now()
			previousMonth := now.AddDate(0, 0, -now.Day())
			stubGenerator.AssertGenerateSingleCalled(t, previousMonth.Year(), int(previousMonth.Month()))
		})

		t.Run("return 400 when only year or month is provided", func(t *testing.T) {
			stubGenerator := report.NewMockGenerator()
			router := testRouterWithReports(stubGenerator)

			missingMonthRecorder := httptest.NewRecorder()
			missingMonthReq, err := http.NewRequest("GET", "/reports/single?year=2022", nil)
			require.NoError(t, err)
			router.ServeHTTP(missingMonthRecorder, missingMonthReq)

			require.Equal(t, http.StatusBadRequest, missingMonthRecorder.Code)
			requireErrorResponse(t, missingMonthRecorder, "either both year and month are provided or none are provided")

			missingYearRecorder := httptest.NewRecorder()
			missingYearReq, err := http.NewRequest("GET", "/reports/single?month=4", nil)
			require.NoError(t, err)
			router.ServeHTTP(missingYearRecorder, missingYearReq)

			require.Equal(t, http.StatusBadRequest, missingYearRecorder.Code)
			requireErrorResponse(t, missingYearRecorder, "either both year and month are provided or none are provided")
		})

		t.Run("return 400 when year is not a number", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/reports/single?year=invalid&month=4", nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			stubGenerator := report.NewMockGenerator()
			router := testRouterWithReports(stubGenerator)

			router.ServeHTTP(recorder, req)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			requireErrorResponse(t, recorder, "year 'invalid' is invalid")
		})

		t.Run("return 400 when year is before 2020", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/reports/single?year=2019&month=4", nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			stubGenerator := report.NewMockGenerator()
			router := testRouterWithReports(stubGenerator)

			router.ServeHTTP(recorder, req)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			requireErrorResponse(t, recorder, "2019 is too early")
		})

		t.Run("return 400 when month is not a number", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/reports/single?year=2022&month=invalid", nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			stubGenerator := report.NewMockGenerator()
			router := testRouterWithReports(stubGenerator)

			router.ServeHTTP(recorder, req)

			require.Equal(t, http.StatusBadRequest, recorder.Code)
			requireErrorResponse(t, recorder, "month 'invalid' is invalid")
		})

		t.Run("return 400 when month is not in range", func(t *testing.T) {
			stubGenerator := report.NewMockGenerator()
			router := testRouterWithReports(stubGenerator)

			for _, m := range []string{"0", "13"} {
				recorder := httptest.NewRecorder()
				req, err := http.NewRequest("GET", "/reports/single?year=2022&month="+m, nil)
				require.NoError(t, err)
				router.ServeHTTP(recorder, req)

				require.Equal(t, http.StatusBadRequest, recorder.Code)
				requireErrorResponse(t, recorder, "month must be in the range 1..12")
			}
		})

		t.Run("return 500 when generator returns error", func(t *testing.T) {
			req, err := http.NewRequest("GET", "/reports/single?year=2022&month=4", nil)
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			stubGenerator := report.NewMockGenerator()
			stubGenerator.StubGenerateSingle().Return(nil, errors.New("some error"))
			router := testRouterWithReports(stubGenerator)

			router.ServeHTTP(recorder, req)

			require.Equal(t, http.StatusInternalServerError, recorder.Code)
			requireErrorResponse(t, recorder, "some error")
		})
	})
}

func requireErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder, message string) {
	var response ErrorResponse
	require.NoError(t, json.NewDecoder(recorder.Body).Decode(&response))
	expectedResponse := ErrorResponse{
		Message: message,
	}
	require.Equal(t, expectedResponse, response)
}

func testRouterWithReports(generator report.Generator) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	RegisterReportsRoutes(r, generator)
	return r
}
