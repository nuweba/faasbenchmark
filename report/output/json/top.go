package json

import (
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"time"
)

const (
	GlobalResultDir = "results"
	LogName         = "faastests.log"
)

type Top struct {
	reportDir string
	logFile   *os.File
}

func New(workingDir string) (report.Top, error) {

	resultDir := filepath.Join(workingDir, GlobalResultDir+"_"+time.Now().Format("20060102150405"))

	t := &Top{reportDir: resultDir}
	err := os.MkdirAll(resultDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, "results dir should be unique")
	}

	return t, nil
}

func (t *Top) LogWriter() (io.Writer, error) {
	logPath := filepath.Join(t.reportDir, LogName)
	f, err := os.Create(logPath)
	if err != nil {
		return nil, err
	}

	t.logFile = f
	return t.logFile, nil
}
