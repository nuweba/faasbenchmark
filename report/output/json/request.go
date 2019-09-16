package json

import (
	"github.com/nuweba/faasbenchmark/report"
)

type Request struct {
	upperLevel *Function
	json       *requestJson
}

type requestJson struct {
	Id                 uint64  `json:"id"`
	InvocationOverHead float64 `json:"invocationOverhead"`
	Duration           float64 `json:"duration"`
	ResponseTime       float64 `json:"responseTime"`
	Reused             bool    `json:"reused"`
	Failed             bool    `json:"failed"`
}

func (f *Function) Request() (report.Request, error) {
	return &Request{upperLevel: f}, nil
}

func (r *Request) Result(result report.Result) error {
	rj := requestJson{
		Id:                 result.Id(),
		InvocationOverHead: result.InvocationOverHead(),
		Duration:           result.Duration(),
		ResponseTime:       result.ContentTransfer(),
		Reused:             result.Reused(),
		Failed:             false,
	}
	r.upperLevel.json.AddResult(rj)
	return nil
}

func (r *Request) Summary(summary string) error {
	return nil
}

func (r *Request) Error(id uint64, error string) error {
	rj := requestJson{
		Id:     id,
		Failed: true,
	}
	r.upperLevel.json.AddResult(rj)
	return nil
}

func (r *Request) RawResult(raw string) error {
	return nil
}
