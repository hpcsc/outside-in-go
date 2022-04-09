package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hpcsc/outside-in-go/internal/report"
	"net/http"
	"strconv"
	"time"
)

const (
	reportsSingleRoutePattern     = "/reports/single"
	reportsCumulativeRoutePattern = "/reports/cumulative"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func RegisterReportsRoutes(router *chi.Mux, generator report.Generator) {
	h := &reportsHandler{
		generator: generator,
	}
	router.Get(reportsSingleRoutePattern, h.Single)
	router.Get(reportsCumulativeRoutePattern, h.Cumulative)
}

type reportsHandler struct {
	generator report.Generator
}

func (h *reportsHandler) Single(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")
	monthParam := r.URL.Query().Get("month")

	year, month, err := h.parseYearAndMonth(yearParam, monthParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.errorResponse(w, err.Error())
		return
	}

	data, err := h.generator.GenerateSingle(*year, *month)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.errorResponse(w, err.Error())
		return
	}

	h.csvResponse(w, year, month, data)
}

func (h *reportsHandler) Cumulative(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")
	monthParam := r.URL.Query().Get("month")

	year, month, err := h.parseYearAndMonth(yearParam, monthParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		h.errorResponse(w, err.Error())
		return
	}

	data, err := h.generator.GenerateCumulative(*year, *month)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.errorResponse(w, err.Error())
		return
	}

	h.csvResponse(w, year, month, data)
}

func (h *reportsHandler) parseYearAndMonth(yearParam string, monthParam string) (*int, *int, error) {
	if (yearParam != "" && monthParam == "") || (yearParam == "" && monthParam != "") {
		return nil, nil, errors.New("either both year and month are provided or none are provided")
	}

	if yearParam == "" && monthParam == "" {
		now := time.Now()
		previousMonth := now.AddDate(0, 0, -now.Day())
		year := previousMonth.Year()
		month := int(previousMonth.Month())
		return &year, &month, nil
	}

	year, err := strconv.Atoi(yearParam)
	if err != nil {
		return nil, nil, fmt.Errorf("year '%s' is invalid", yearParam)
	}

	if year < 2020 {
		return nil, nil, fmt.Errorf("%d is too early", year)
	}

	month, err := strconv.Atoi(monthParam)
	if err != nil {
		return nil, nil, fmt.Errorf("month '%s' is invalid", monthParam)
	}

	if month < 1 || month > 12 {
		return nil, nil, errors.New("month must be in the range 1..12")
	}

	return &year, &month, nil
}

func (h *reportsHandler) errorResponse(w http.ResponseWriter, message string) error {
	return json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
	})
}

func (h *reportsHandler) csvResponse(w http.ResponseWriter, year *int, month *int, data []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=single-%d%02d.csv", *year, *month))
	w.Write(data)
}
