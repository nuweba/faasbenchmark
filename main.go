package main

import (
	"fmt"
	"github.com/nuweba/faasbenchmark/cmd"
	"os/exec"
)

func isInstalled(bin, name string) bool {
	_, err := exec.LookPath(bin)
	if err != nil {
		fmt.Printf("%s is not installed\n", name)
		return false
	}
	return true
}

func checkDeps() bool {
	dependencies := map[string]string{
		"sls":    "serverless framework",
		"az":     "azure-cli",
		"dotnet": "dotnet-sdk",
		"mvn":    "maven",
		"func":   "azure-functions-core-tools"}

	depsInstalled := true
	for binary, name := range dependencies {
		depsInstalled = isInstalled(binary, name) && depsInstalled
	}
	return depsInstalled
}

func main() {
	if !checkDeps() {
		fmt.Println("please install all missing dependencies")
		return
	}
	cmd.Execute()
}
