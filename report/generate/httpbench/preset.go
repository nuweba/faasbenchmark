package httpbench

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/httpbench"
)


func ReportFunctionResults(funcConfig *config.HttpFunction, rr *httpbench.PresetResult) (string, error) {
	funcConfig.Logger.Debug("marshaling function result")
	jsonResult, err := yaml.Marshal(rr)
	if err != nil {
		return "", err
	}

	funcConfig.Logger.Debug("writing function result")
	functionReportW, err := funcConfig.Report.ResultWriter()
	if err != nil {
		return "", errors.Wrap(err, "function result")
	}
	fmt.Fprintln(functionReportW, string(jsonResult))
	funcConfig.Logger.Info("finished function test", zap.String("result", string(jsonResult)))
	return string(jsonResult), nil
}
