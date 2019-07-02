package multi

import (
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
)

type multiRequest struct {
	request []report.Request
}

func (mReq *multiRequest) Result(result report.Result) error {
	for _, t := range mReq.request {
		err := t.Result(result)
		if err != nil {
			return errors.Wrap(err, "multi request result writer error")
		}
	}

	return nil
}

func (mReq *multiRequest) Summary(summary string) error {
	for _, t := range mReq.request {
		err := t.Summary(summary)
		if err != nil {
			return errors.Wrap(err, "multi request summary writer error")
		}
	}

	return nil
}

func (mReq *multiRequest) Error(id uint64, error string) error {
	for _, t := range mReq.request {
		err := t.Error(id, error)
		if err != nil {
			return errors.Wrap(err, "multi request error writer error")
		}
	}

	return nil
}

func (mReq *multiRequest) RawResult(raw string) error {
	for _, t := range mReq.request {
		err := t.RawResult(raw)
		if err != nil {
			return errors.Wrap(err, "multi request raw result writer error")
		}
	}

	return nil
}
