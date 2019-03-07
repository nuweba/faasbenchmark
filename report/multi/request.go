package multi

import (
	"github.com/pkg/errors"
	"io"
	"github.com/nuweba/faasbenchmark/report"
)

type multiRequest struct {
	request []report.Request
}

func (mReq *multiRequest) ResultWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mReq.request {
		writer, err := t.ResultWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi request result writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}

func (mReq *multiRequest) SummaryWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mReq.request {
		writer, err := t.SummaryWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi request summary writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}


func (mReq *multiRequest) ErrorWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mReq.request {
		writer, err := t.ErrorWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi request error writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}

func (mReq *multiRequest) RawResultWriter() (io.Writer, error) {
	var writers []io.Writer
	for _, t := range mReq.request {
		writer, err := t.RawResultWriter()
		if err != nil {
			return nil, errors.Wrap(err, "multi request raw result writer error")
		}
		writers = append(writers, writer)
	}

	return io.MultiWriter(writers...), nil
}