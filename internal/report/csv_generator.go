package report

import "github.com/hpcsc/outside-in-go/internal/storer"

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
	existingAggregated, _ := g.storer.RetrieveAggregated(storer.SingleReportType, year, month)
	if existingAggregated != nil {
		return existingAggregated, nil
	}

	return nil, nil
}

func (g *csvGenerator) GenerateCumulative(year int, month int) ([]byte, error) {
	return nil, nil
}
