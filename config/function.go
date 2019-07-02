package config

import (
	"github.com/nuweba/faasbenchmark/report"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
)

type HttpFunction struct {
	Test       *Test
	HttpConfig *Http
	Function   stack.Function
	Report     report.Function
	Logger    *zap.Logger
}

func (hf *HttpFunction) newLogger(writer io.Writer) *zap.Logger {
	output := zapcore.Lock(zapcore.AddSync(writer))

	cfg := zap.NewDevelopmentEncoderConfig()
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), output, zap.DebugLevel)
	//l := zap.New(core, zap.Option(zap.Development()), zap.Option(zap.AddCaller()))

	globalLogger := hf.Test.Config.Logger.With(zap.Namespace(hf.Function.Name()))
	hf.Logger = zap.New(zapcore.NewTee(core, globalLogger.Core()), zap.Option(zap.Development()), zap.Option(zap.AddCaller()))
	hf.Logger = hf.Logger.Named("test").Named(hf.Test.TestId)
	return hf.Logger

}

func (t *Test) NewFunction(httpConfig *Http, function stack.Function) (*HttpFunction, error) {
	functionReport, err := t.Report.Function(function.Name(), function.Description(), function.Runtime(), function.MemorySize())

	if err != nil {
		t.Config.Logger.Error("function report", zap.Error(err))
		return nil, errors.Wrap(err, "function Report")
	}

	hf := &HttpFunction{
		Test:       t,
		HttpConfig: httpConfig,
		Function:   function,
		Report:     functionReport,
	}

	loggerW, err := functionReport.LogWriter()
	if err != nil {
		t.Config.Logger.Error("log writer", zap.Error(err))
		return nil, err
	}

	l := hf.newLogger(loggerW)

	l.Debug("testing function", zap.String("name", function.Name()), zap.String("description", function.Description()))

	err = hf.Report.StackDescription(t.Stack.Description)
	if err != nil {
		t.Config.Logger.Error("description writer", zap.Error(err))
		return nil, errors.Wrap(err, "function report description writer")
	}

	httpConfigRaw, err := httpConfig.String()
	if err!= nil {
		t.Config.Logger.Error("HttpTestConfig", zap.Error(err))
		return nil, errors.Wrap(err, "function report http test config to string")
	}

	err = hf.Report.HttpTestConfig(httpConfigRaw)
	if err != nil {
		t.Config.Logger.Error("HttpTestConfig", zap.Error(err))
		return nil, errors.Wrap(err, "function report http test config writer")
	}

	return hf, nil
}
