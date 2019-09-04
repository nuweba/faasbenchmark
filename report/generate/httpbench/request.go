package httpbench

import (
	"encoding/json"
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/report"
	"github.com/nuweba/httpbench/engine"
	"github.com/nuweba/httpbench/syncedtrace"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"net/http"
	"time"
)

const MaxBodySize = 512

type IsWarm bool
type FunctionOutput struct {
	Reused   bool          `json:"reused"`
	Duration time.Duration `json:"duration"`
}

func Reuse(s bool) string {
	switch s {
	case true:
		return "reused"
	default:
		return "fresh"
	}
}

func RequestBodyUnmarshal(body []byte) (*FunctionOutput, error) {
	funcOutput := &FunctionOutput{}
	err := json.Unmarshal(body, funcOutput)
	if err != nil {
		return nil, err
	}

	if funcOutput.Duration == 0 {
		return nil, errors.New(fmt.Sprint("duration field is missing", string(body)))
	}

	return funcOutput, nil
}

func TraceResultString(tr engine.TraceResult) (string, error) {
	//truncating body if too large to save log space
	if len(tr.Body) > MaxBodySize {
		tr.Body = truncate(tr.Body)
	}
	if tr.Err != nil && len(tr.Err.Message) > MaxBodySize {
		tr.Err.Message = truncate(tr.Err.Message)
	}

	b, err := yaml.Marshal(tr)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func truncate(s string) string {
	resultStart := s[0 : MaxBodySize/2]
	resultEnd := s[len(s)-(MaxBodySize/2) : len(s)-1]
	return fmt.Sprintf("%s...truncated...%s", resultStart, resultEnd)
}

func TraceResultSummaryError(httpConf *config.Http, tr *engine.TraceResult, funcOutput *FunctionOutput) *Summary {
	summaryWithErrors := TraceResultSummaryHttps(httpConf, tr, funcOutput)
	summaryWithErrors.Errors = make(map[string]error)
	for hook := syncedtrace.TraceHookType(0); hook < syncedtrace.HooksCount; hook++ {
		if tr.Hooks[hook].Err != nil {
			summaryWithErrors.Errors[hook.String()] = tr.Hooks[hook].Err
		}
	}
	return summaryWithErrors
}

type Summary struct {
	Id     uint64
	Status string
	Start  time.Time
	*engine.TraceSummary
	FunctionDuration time.Duration
	Done             time.Time
	Total            time.Duration
	Errors           map[string]error
}

func (s *Summary) String() string {
	b, err := yaml.Marshal(s)
	if err != nil {
		return ""
	}

	return string(b)
}

func TraceResultSummaryHttps(httpConf *config.Http, tr *engine.TraceResult, funcOutput *FunctionOutput) *Summary {
	return &Summary{
		Id:               tr.Id,
		Status:           Reuse(funcOutput.Reused),
		Start:            tr.Start,
		TraceSummary:     tr.Summary,
		FunctionDuration: funcOutput.Duration - httpConf.SleepTime,
		Done:             tr.Done,
		Total:            tr.Total - httpConf.SleepTime,
	}
}

type errorReport struct {
	reporter report.Request
	logger   *zap.Logger
}

func (e *errorReport) errorReporter(id uint64, err error, errStr string, data string) {
	e.logger.Error(errStr, zap.Error(err), zap.String("summary", data))

	err2 := e.reporter.Error(id, fmt.Sprintf("Error id: %d, msg:%s, error: %s, data:%s\n", id, errStr, err.Error(), data))
	if err2 != nil {
		e.logger.Error("report error writer", zap.Error(err2))
	}
}

func ReportRequestResults(funcConfig *config.HttpFunction, resultCh chan *engine.TraceResult, outputFn provider.RequestFilter, debug bool) {
	reqReport, err := funcConfig.Report.Request()
	if err != nil {
		funcConfig.Logger.Fatal("request report", zap.Error(err))
		return
	}

	errorReporter := &errorReport{reporter: reqReport, logger: funcConfig.Logger}

	for result := range resultCh {
		raw, err := TraceResultString(*result)
		if err != nil {
			funcConfig.Logger.Error("error marshaling trace result", zap.Error(err))
		} else {
			funcConfig.Logger.Debug("writing raw result")
			err = reqReport.RawResult(raw)
			if err != nil {
				funcConfig.Logger.Error("raw result writer", zap.Error(err))
			}
		}

		funcConfig.Logger.Debug("got new request result", zap.Uint64("id", result.Id))
		funcOutput, err := RequestBodyUnmarshal([]byte(result.Body))
		if err != nil {
			errorReporter.errorReporter(result.Id, err, "request body unmarshal", result.Body)
			continue
		}

		if result.Err != nil {
			errorReporter.errorReporter(result.Id, result.Err, "trace error", TraceResultSummaryError(funcConfig.HttpConfig, result, funcOutput).String())
			continue
		}

		if result.Error {
			errorReporter.errorReporter(result.Id, errors.New("result error without error msg"), "trace error", TraceResultSummaryError(funcConfig.HttpConfig, result, funcOutput).String())
			continue
		}

		if result.Response.StatusCode != http.StatusOK {
			errorReporter.errorReporter(result.Id, result.Err, "function did not return 200 ok", TraceResultSummaryError(funcConfig.HttpConfig, result, funcOutput).String())
			continue
		}

		summary := TraceResultSummaryHttps(funcConfig.HttpConfig, result, funcOutput)
		funcConfig.Logger.Debug("trace result", zap.Any("summary", summary))

		funcConfig.Logger.Debug("running filter function on result")
		filteredResult, err := outputFn(funcConfig.HttpConfig.SleepTime, result, funcOutput.Duration, funcOutput.Reused)
		if err != nil {
			errorReporter.errorReporter(result.Id, result.Err, "output fn", "")
			continue
		}

		funcConfig.Logger.Debug("filter function result",
			zap.Uint64("id", filteredResult.Id()),
			zap.Float64("invocationoverhead", filteredResult.InvocationOverHead()),
			zap.Float64("duration", filteredResult.Duration()),
			zap.Float64("contenttransfer", filteredResult.ContentTransfer()),
			zap.Bool("reused", filteredResult.Reused()),
		)
		err = reqReport.Result(filteredResult)
		if err != nil {
			funcConfig.Logger.Error("result writer", zap.Error(err))
		}

		summaryOutput := fmt.Sprintf("%d: InvocationOverHead: %f, duration: %f\n", result.Id, filteredResult.InvocationOverHead(), filteredResult.Duration())
		if debug {
			summaryOutput += summary.String() + "\n"
		}
		err = reqReport.Summary(summaryOutput)
		if err != nil {
			funcConfig.Logger.Error("summary writer", zap.Error(err))
		}

		funcConfig.Logger.Info("request done", zap.Uint64("id", result.Id))
	}

}

func DiscardRequestResults(resultCh chan *engine.TraceResult) {
	for range resultCh {
	}
}
