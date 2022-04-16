package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/hpcsc/outside-in-go/internal/handler"
	"github.com/hpcsc/outside-in-go/internal/report"
	"github.com/hpcsc/outside-in-go/internal/storer"
	"log"
	"net/http"
	"os"
)

var (
	s3Endpoint = os.Getenv("S3_ENDPOINT")
	port       = os.Getenv("PORT")
	bucket     = os.Getenv("BUCKET")
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	handler.RegisterReportsRoutes(r, newReportGenerator())

	addr := fmt.Sprintf(":%s", port)
	log.Printf("listening at %s", addr)
	http.ListenAndServe(addr, r)
}

func newReportGenerator() report.Generator {
	s, err := storer.NewS3Storer(s3Endpoint, bucket)
	if err != nil {
		log.Fatalf("%v", err)
	}
	gen := report.NewCsvGenerator(s)
	return gen
}
