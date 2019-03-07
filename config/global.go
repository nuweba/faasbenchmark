package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/report"
	"nwb.nu/httpbench/engine"
)

type Global struct {
	Provider  provider.FaasProvider
	Stacks    *Stacks
	report    report.Top
	resultDir string
	Logger    *zap.Logger
	logCh     chan *engine.TraceResult
}

func newLogger(writer io.Writer) *zap.Logger {
	output := zapcore.Lock(zapcore.AddSync(writer))

	cfg := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), output, zap.DebugLevel)
	l := zap.New(core, zap.Option(zap.Development()), zap.Option(zap.AddCaller()))

	return l.Named("main")

}

func NewGlobalConfig(provider provider.FaasProvider, arsenalPath string, report report.Top) (*Global, error) {

	loggerW, err := report.LogWriter()
	if err != nil {
		return nil, err
	}

	l := newLogger(loggerW)

	l.Info("starting tests")

	l = l.Named(provider.Name())

	stacks, err := newStacks(provider, arsenalPath)
	if err != nil {
		return nil, err
	}

	l.Debug("stacks loaded", zap.String("arsenal", arsenalPath))
	return &Global{report: report, Logger: l, Provider: provider, Stacks: stacks}, nil
}
