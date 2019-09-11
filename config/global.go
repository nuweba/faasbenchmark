package config

import (
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/report"
	"github.com/nuweba/httpbench/engine"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type Global struct {
	Provider  provider.FaasProvider
	Stacks    *Stacks
	report    report.Top
	resultDir string
	Logger    *zap.Logger
	logCh     chan *engine.TraceResult
	Debug     bool
}

func newLogger(writer io.Writer, debug bool) *zap.Logger {
	lvl := zap.DebugLevel
	if !debug {
		lvl = zap.InfoLevel
	}
	output := zapcore.Lock(zapcore.AddSync(writer))

	cfg := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), output, lvl)
	l := zap.New(core, zap.Option(zap.Development()), zap.Option(zap.AddCaller()))

	return l.Named("main")

}

func NewGlobalConfig(provider provider.FaasProvider, arsenalPath string, report report.Top, debug bool) (*Global, error) {

	loggerW, err := report.LogWriter()
	if err != nil {
		return nil, err
	}

	l := newLogger(loggerW, debug)

	l.Info("starting tests")

	l = l.Named(provider.Name())

	stacks, err := newStacks(provider, arsenalPath)
	if err != nil {
		return nil, err
	}

	l.Debug("stacks loaded", zap.String("arsenal", arsenalPath))
	return &Global{report: report, Logger: l, Provider: provider, Stacks: stacks, Debug: debug}, nil
}
