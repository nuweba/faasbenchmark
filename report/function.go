package report

import (
	"io"
)

type Function interface {
	LogWriter() (io.Writer, error)
	ResultWriter() (io.Writer, error)
	DescriptionWriter() (io.Writer, error)
	HttpTestConfigWriter() (io.Writer, error)
	Request() (Request, error)
}
