package file

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

const (
	RawResultDir = "raw"
	SummaryPath  = "full.log"
	ErrorPath    = "error.log"
)

type Request struct {
	upperLevel            *Function
	functionReqResultFile *os.File
	rawResultDir          string
	rawResultFile         *os.File
	SummaryFile           *os.File
	ErrorFile             *os.File
}

func (f *Function) Request() (report.Request, error) {
	r := &Request{upperLevel: f}

	//result, just x,y
	functionReqResultFilePath := filepath.Join(r.upperLevel.functionResultPath, r.upperLevel.functionName)
	functionReqResultFile, err := os.OpenFile(functionReqResultFilePath, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function result file should be unique")
	}
	r.functionReqResultFile = functionReqResultFile

	//full summery
	functionSummaryFilePath := filepath.Join(r.upperLevel.functionResultPath, r.upperLevel.functionName+"_"+SummaryPath)
	summaryFile, err := os.OpenFile(functionSummaryFilePath, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function error file should be unique")
	}
	r.SummaryFile = summaryFile

	//error file
	functionErrorFilePath := filepath.Join(r.upperLevel.functionResultPath, r.upperLevel.functionName+"_"+ErrorPath)
	functionErrorFile, err := os.OpenFile(functionErrorFilePath, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "function summary file should be unique")
	}
	r.ErrorFile = functionErrorFile

	//raw result file and dir
	rawResultDir := filepath.Join(r.upperLevel.functionResultPath, RawResultDir)
	err = os.MkdirAll(rawResultDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "raw request result dir should be unique")
	}

	functionRawResultFilePath := filepath.Join(rawResultDir, r.upperLevel.functionName)
	functionRawResultFile, err := os.OpenFile(functionRawResultFilePath, os.O_APPEND|os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "raw request result file should be unique")
	}
	r.rawResultDir = rawResultDir
	r.rawResultFile = functionRawResultFile

	return r, nil
}

func (r *Request) Result(result report.Result) error {
	t := fmt.Sprintf("id: %s, invocationOverHead: %s, duration: %s, responseTime: %s, reused: %s\n", result.Id(), result.InvocationOverHead(), result.Duration(), result.ContentTransfer(), result.Reused())
	_, err := r.functionReqResultFile.WriteString(t)
	return err
}

func (r *Request) Summary(summary string) error {
	_, err := r.SummaryFile.WriteString(summary)
	return err
}

func (r *Request) Error(id uint64, error string) error {
	_, err := r.ErrorFile.WriteString(error)
	return err
}

func (r *Request) RawResult(raw string) error {
	_, err := r.rawResultFile.WriteString(raw)
	return err
}
