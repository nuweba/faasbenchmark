package json

import (
	"github.com/nuweba/faasbenchmark/report"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

const (
	GlobalResultDir = "results"
)

type Top struct {
	reportDir string
	logFile   *os.File
}

func New(workingDir string) (report.Top, error) {

	resultDir := filepath.Join(workingDir, GlobalResultDir+"_"+time.Now().Format("20060102150405"))

	t := &Top{reportDir: resultDir}
	err := os.MkdirAll(resultDir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, errors.Wrap(err, "could not create results dir")
	}

	return t, nil
}

func (t *Top) LogWriter() (io.Writer, error) {
	return ioutil.Discard, nil
}
