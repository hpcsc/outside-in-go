package storer

type ReportType string

const (
	SingleReportType     ReportType = "single"
	CumulativeReportType ReportType = "cumulative"
)

type Storer interface {
	RetrieveIndividualFiles(reportType ReportType, year int, month int) ([][]byte, error)
	RetrieveAggregated(reportType ReportType, year int, month int) ([]byte, error)
	StoreAggregated(reportType ReportType, year int, month int, data []byte) error
}
