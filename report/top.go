package report

import (
	"io"
)

type Top interface {
	LogWriter() (io.Writer, error)
	Test(testId string, providerName string) (Test, error)
}
