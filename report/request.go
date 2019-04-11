package report

type Request interface {
	Result(result string) error
	Summary(summary string) error
	Error(error string) error
	RawResult(raw string) error
}
