package file

import (
	"github.com/pkg/errors"
	"io"
	"github.com/nuweba/faasbenchmark/report"
	"os"
	"path/filepath"
)

const (
	RawResultDir = "raw"
	SummaryPath  = "full.log"
	ErrorPath  = "error.log"
)

type Request struct {
	upperLevel *Function
	functionReqResultFile *os.File
	rawResultDir string
	rawResultFile *os.File
	SummaryFile *os.File
	ErrorFile *os.File
}

func (f *Function) Request() (report.Request, error) {
	r := &Request{upperLevel: f}

	return r, nil
}

func (r *Request) ResultWriter() (io.Writer, error) {
	functionReqResultFile := filepath.Join(r.upperLevel.functionResultPath, r.upperLevel.functionName)
	file, err := os.OpenFile(functionReqResultFile,  os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function result file should be unique")
	}
	r.functionReqResultFile = file
	return r.functionReqResultFile, err
}

func (r *Request) SummaryWriter() (io.Writer, error) {
	functionSummaryFile := filepath.Join(r.upperLevel.functionResultPath, r.upperLevel.functionName + "_" +SummaryPath)
	file, err := os.OpenFile(functionSummaryFile,  os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function error file should be unique")
	}
	r.SummaryFile = file
	return r.SummaryFile, err
}

func (r *Request) ErrorWriter() (io.Writer, error) {
	functionErrorFile := filepath.Join(r.upperLevel.functionResultPath, r.upperLevel.functionName + "_" +ErrorPath)
	file, err := os.OpenFile(functionErrorFile,  os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function summary file should be unique")
	}
	r.ErrorFile = file
	return r.ErrorFile, err
}

func (r *Request) RawResultWriter() (io.Writer, error) {
	rawResultDir := filepath.Join(r.upperLevel.functionResultPath, RawResultDir)
	err := os.MkdirAll(rawResultDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "raw request result dir should be unique")
	}

	functionReqResultFile := filepath.Join(rawResultDir, r.upperLevel.functionName)
	file, err := os.OpenFile(functionReqResultFile,  os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "raw request result file should be unique")
	}
	r.rawResultDir = rawResultDir
	r.rawResultFile = file
	return r.rawResultFile, err
}
