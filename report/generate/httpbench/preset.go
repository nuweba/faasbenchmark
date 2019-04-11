package httpbench

import (
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/httpbench"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

func ReportFunctionResults(funcConfig *config.HttpFunction, rr *httpbench.PresetResult) (string, error) {
	funcConfig.Logger.Debug("marshaling function result")
	jsonResult, err := yaml.Marshal(rr)
	if err != nil {
		return "", err
	}

	funcConfig.Logger.Debug("writing function result")
	err = funcConfig.Report.BenchResult(string(jsonResult))
	if err != nil {
		return "", errors.Wrap(err, "function result")
	}
	funcConfig.Logger.Info("finished function test", zap.String("result", string(jsonResult)))
	return string(jsonResult), nil
}
