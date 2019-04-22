package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/nuweba/faasbenchmark/provider"
	"github.com/nuweba/faasbenchmark/testsuite"
	"strings"
)

func init() {
	var echoTimes int

	var cmdList = &cobra.Command{
		Use:   "list [COMMAND]",
		Short: "list the providers or tests",
		Long: `list will show all the current providers and test functions`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not supported: " + strings.Join(args, " "))
		},
	}

	var cmdListTests = &cobra.Command{
		Use:   "tests",
		Short: "show all the test id's",
		Args: cobra.NoArgs,
		Run:  listTests,
	}

	var cmdListProviders = &cobra.Command{
		Use:   "providers",
		Short: "show all the providers",
		Args: cobra.NoArgs,
		Run:  listProviders,
	}

	cmdListTests.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	rootCmd.AddCommand(cmdList)
	cmdList.AddCommand(cmdListTests, cmdListProviders)
}

func listTests(cmd *cobra.Command, args []string) {
	_ = cmd
	_ = args
	//cloudTests, err := provider.New("aws")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	counter := 1
	for id, test := range testsuite.Tests.TestFunctions {
		fmt.Println(counter, id, "-", test.Description)
		counter++
	}
}

func listProviders(cmd *cobra.Command, args []string) {
	providers := provider.ProviderList()
	for i, provider := range providers {
		fmt.Println(int(i) + 1, provider)
	}
}
