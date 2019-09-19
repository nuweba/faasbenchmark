package report

type Result interface {
	Id() uint64
	InvocationOverHead() float64
	Duration() float64
	ContentTransfer() float64
	Reused() bool
}

type Request interface {
	Result(result Result) error
	Summary(summary string) error
	Error(id uint64, error string) error
	RawResult(raw string) error
}
