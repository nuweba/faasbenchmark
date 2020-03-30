package cmd

import (
	"errors"
	"fmt"
	"github.com/nuweba/faasbenchmark/config"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/report"
	"github.com/nuweba/faasbenchmark/report/multi"
	"github.com/nuweba/faasbenchmark/report/output/file"
	"github.com/nuweba/faasbenchmark/report/output/json"
	"github.com/nuweba/faasbenchmark/report/output/stdio"
	"github.com/nuweba/faasbenchmark/testsuite"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

const (
	TestsDir           = "arsenal"
	exampleTestsPrefix = "example"
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
		cmdProvider.Flags().BoolVarP(&debug, "debug", "d", false, "whether to show debug output, default is false")
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
	var report report.Top
	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		return errors.New("can't get pkg dir")
	}

	pgkPath := filepath.Join(path.Dir(filename), "/../") // ./pkg/cmd/.../

	jsonReport, err := json.New(resultPath)
	if err != nil {
		return err
	}

	stdioReport, err := stdio.New(os.Stdout)
	if err != nil {
		return err
	}

	if debug {
		fileReport, err := file.New(resultPath)
		if err != nil {
			return err
		}
		report = multi.Report(fileReport, jsonReport, stdioReport)
	} else {
		report = multi.Report(jsonReport, stdioReport)
	}

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

func RunAllTests(gConfig *config.Global) (err error) {
	for id := range testsuite.Tests.TestFunctions {
		if strings.HasPrefix(strings.ToLower(id), exampleTestsPrefix) {
			continue
		}
		testErr := runOneTest(gConfig, id)
		if testErr != nil {
			gConfig.Logger.Error("error running test", zap.Error(err), zap.String("test", id))
			err = testErr
		}
	}
	return err
}

func RunSpecificTests(gConfig *config.Global, testIds ...string) (err error) {
	for _, id := range testIds {
		testErr := runOneTest(gConfig, id)
		if testErr != nil {
			gConfig.Logger.Error("error running test", zap.Error(err), zap.String("test", id))
			err = testErr
		}
	}
	return err
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
	stackRemoved := make(chan struct{})
	go handleSignals(gConfig, stack, stackRemoved)
	defer func() {
		if err != nil {
			// error deploying stack, will not remove it
			return
		}
		gConfig.Logger.Info("removing stack", zap.String("name", stack.StackId()))
		// err will be returned by wrapping function's return statement
		removeErr := stack.RemoveStack()
		if removeErr != nil {
			err = removeErr
			gConfig.Logger.Warn("failed removing stack", zap.String("name", stack.StackId()), zap.String("err", err.Error()))
		} else {
			stackRemoved <- struct{}{}
			gConfig.Logger.Debug("stack removed", zap.String("name", stack.StackId()))
		}
	}()
	err = stack.DeployStack()
	if err != nil {
		return err
	}
	gConfig.Logger.Debug("stack deployed", zap.String("name", stack.StackId()))

	gConfig.Logger.Info("running test", zap.String("name", test.Id))

	test.Fn(testConfig)

	gConfig.Logger.Debug("test is done", zap.String("name", test.Id))

	return err
}

func handleSignals(gConfig *config.Global, stack *config.Stack, stackRemoved chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-signals:
		err := stack.RemoveStack()
		if err != nil {
			gConfig.Logger.Error("removing stack failed", zap.String("err", err.Error()))
		}
		os.Exit(1)
	case <-stackRemoved:
		return
	}
}
