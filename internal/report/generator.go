package report

type Generator interface {
	GenerateSingle(year int, month int) ([]byte, error)
	GenerateCumulative(year int, month int) ([]byte, error)
}
