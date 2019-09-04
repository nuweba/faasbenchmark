package cmd

import (
	"errors"
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/report/multi"
	"github.com/nuweba/faasbenchmark/report/output/json"
	"github.com/nuweba/faasbenchmark/report/output/stdio"
	"github.com/nuweba/faasbenchmark/testsuite"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	TestsDir = "arsenal"
)

var resultPath string
var debug bool

func init() {

	var cmdRun = &cobra.Command{
		Use:   "run [provider] [test id]",
		Short: "run a specific test",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("provider not supported: " + strings.Join(args, " "))
		},
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	cmdRun.Flags().StringVarP(&resultPath, "resultPath", "r", dir, "directory to write the results, default is cwd")
	cmdRun.Flags().BoolVarP(&debug, "debug", "d", false, "whether to show debug output, default is false")
	rootCmd.AddCommand(cmdRun)

	for providerId := provider.Providers(0); providerId < provider.ProvidersCount; providerId++ {
		p := providerId
		cmdProvider := &cobra.Command{
			Use:   fmt.Sprintf("%s [test id / all]", p.String()),
			Short: p.Description(),
			Args:  validateTestName,
			Run: func(cmd *cobra.Command, args []string) {

				err := runTests(p.String(), args...)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

			},
		}
		cmdRun.AddCommand(cmdProvider)
	}
}

func validateTestName(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("requires at least one test id or \"all\"")
	}

	if args[0] == "all" {
		return nil
	}

	if _, err := testsuite.Tests.GetTestSuite(args[0]); err != nil {
		return err
	}

	return nil
}

func runTests(providerName string, testIds ...string) error {
	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		return errors.New("can't get pkg dir")
	}

	pgkPath := filepath.Join(path.Dir(filename), "/../") // ./pkg/cmd/.../

	//todo: changeme
	//fileReport, err := file.New(resultPath)
	fileReport, err := json.New(resultPath)
	if err != nil {
		return err
	}

	stdioReport, err := stdio.New(os.Stdout)

	if err != nil {
		return err
	}

	report := multi.Report(fileReport, stdioReport)

	faasProvider, err := provider.NewProvider(providerName)
	if err != nil {
		return err
	}

	arsenalPath := filepath.Join(pgkPath, TestsDir)
	gConfig, err := config.NewGlobalConfig(faasProvider, arsenalPath, report, debug)
	if err != nil {
		return err
	}

	switch testIds[0] {
	case "all":
		err = RunAllTests(gConfig)
	default:
		err = RunSpecificTests(gConfig, testIds...)
	}

	if err != nil {
		gConfig.Logger.Error("error running test", zap.Error(err))
		return errors.New(fmt.Sprint("error in run Tests:", err))

	}

	return nil
}

func RunAllTests(gConfig *config.Global) error {
	for id := range testsuite.Tests.TestFunctions {
		err := runOneTest(gConfig, id)
		if err != nil {
			gConfig.Logger.Error("error running test", zap.Error(err), zap.String("test", id))
			continue
		}
	}

	return nil
}

func RunSpecificTests(gConfig *config.Global, testIds ...string) error {
	for _, id := range testIds {
		err := runOneTest(gConfig, id)
		if err != nil {
			gConfig.Logger.Error("error running test", zap.Error(err), zap.String("test", id))
			continue
		}
	}

	return nil
}

func runOneTest(gConfig *config.Global, testId string) error {
	test, err := testsuite.Tests.GetTestSuite(testId)
	if err != nil {
		return err
	}

	gConfig.Logger.Info("got test suite", zap.String("name", test.Id), zap.String("description", test.Description))

	stack, err := gConfig.Stacks.GetStack(test.RequiredStack)
	if err != nil {
		return err
	}
	gConfig.Logger.Info("got stack", zap.String("name", stack.StackId()))
	testConfig, err := gConfig.NewTest(stack, test.Id, test.Description)

	gConfig.Logger.Info("deploying stack", zap.String("name", stack.StackId()))
	err = stack.DeployStack()
	if err != nil {
		return err
	}
	gConfig.Logger.Debug("stack deployed", zap.String("name", stack.StackId()))

	gConfig.Logger.Info("running test", zap.String("name", test.Id))

	test.Fn(testConfig)

	gConfig.Logger.Debug("test is done", zap.String("name", test.Id))

	gConfig.Logger.Info("removing stack", zap.String("name", stack.StackId()))
	err = stack.RemoveStack()
	if err != nil {
		return err
	}
	gConfig.Logger.Debug("stack removed", zap.String("name", stack.StackId()))

	return nil
}
