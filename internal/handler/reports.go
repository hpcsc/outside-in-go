package handler

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/hpcsc/outside-in-go/internal/report"
	"net/http"
	"strconv"
)

const (
	reportsSingleRoutePattern = "/reports/single"
)

func RegisterReportsRoutes(router *chi.Mux, generator report.Generator) {
	h := &reportsHandler{
		generator: generator,
	}
	router.Get(reportsSingleRoutePattern, h.Single)
}

type reportsHandler struct {
	generator report.Generator
}

func (h *reportsHandler) Single(w http.ResponseWriter, r *http.Request) {
	yearParam := r.URL.Query().Get("year")
	monthParam := r.URL.Query().Get("month")

	year, _ := strconv.Atoi(yearParam)
	month, _ := strconv.Atoi(monthParam)

	data, _ := h.generator.GenerateSingle(year, month)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=single-%d%02d.csv", year, month))
	w.Write(data)
}
