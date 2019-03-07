package report

import (
	"io"
)

type Request interface {
	ResultWriter() (io.Writer, error)
	SummaryWriter() (io.Writer, error)
	ErrorWriter() (io.Writer, error)
	RawResultWriter() (io.Writer, error)
}
