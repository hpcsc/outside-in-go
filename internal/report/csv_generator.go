package report

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/hpcsc/outside-in-go/internal/storer"
	"log"
)

var _ Generator = &csvGenerator{}

func NewCsvGenerator(storer storer.Storer) Generator {
	return &csvGenerator{
		storer: storer,
	}
}

type csvGenerator struct {
	storer storer.Storer
}

func (g *csvGenerator) GenerateSingle(year int, month int) ([]byte, error) {
	existingAggregated, err := g.storer.RetrieveAggregated(storer.SingleReportType, year, month)
	if err != nil {
		log.Printf("failed to retrieved existing aggregate: %v", err)
		// continue with aggregate logic
	} else if existingAggregated != nil {
		return existingAggregated, nil
	}

	files, err := g.storer.RetrieveIndividualFiles(storer.SingleReportType, year, month)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no data available for %02d/%d", month, year)
	}

	var aggregatedLines [][]string
	for i, f := range files {
		reader := csv.NewReader(bytes.NewReader(f))
		lines, err := reader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("failed to read one or more csv files: %v", err)
		}

		if i == 0 {
			aggregatedLines = append(aggregatedLines, lines[0])
		}

		aggregatedLines = append(aggregatedLines, lines[1:]...)
	}

	var aggregated bytes.Buffer
	writer := csv.NewWriter(&aggregated)
	if err = writer.WriteAll(aggregatedLines); err != nil {
		return nil, fmt.Errorf("failed to write csv content to buffer: %v", err)
	}

	if err := g.storer.StoreAggregated(storer.SingleReportType, year, month, aggregated.Bytes()); err != nil {
		return nil, err
	}

	return aggregated.Bytes(), nil
}

func (g *csvGenerator) GenerateCumulative(year int, month int) ([]byte, error) {
	return nil, nil
}
