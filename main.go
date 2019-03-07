package main

import "github.com/nuweba/faasbenchmark/cmd"

func main() {
	cmd.Execute()
	//req, err := http.NewRequest("GET", "https://www.google.co.il", bytes.NewReader([]byte("")))
	//if err != nil {
	//	fmt.Println("Error in new request")
	//}
	//duration := 5 * time.Second
	//reqDelay := 50 * time.Millisecond
	//concurrencyLimit := uint64(10)
	//httptrace.RequestPerDuration(req, reqDelay, duration)
	//httptrace.ConcurrentRequestsSynced(req,concurrencyLimit, reqDelay, duration)
	//httptrace.ConcurrentRequestsUnsynced(req,concurrencyLimit, reqDelay, duration)


	//hitsGraph := httptrace.HitsGraph(
	//	[]httptrace.RequestsPerTime{
	//		{2, 1 * time.Second},
	//		{3, 1 * time.Second},
	//		{4, 1 * time.Second},
	//	},
	//	)
	//
	//httptrace.RequestsForTimeGraph(req, hitsGraph)


	//concurrentGraph := httptrace.ConcurrentGraph(
	//	[]httptrace.RequestsPerTime{
	//		{2, 0 * time.Second},
	//		{1, 1 * time.Second},
	//		{8, 10 * time.Second},
	//		{2, 13 * time.Second},
	//		{1, 20 * time.Second},
	//	},
	//	)
	//
	//httptrace.ConcurrentForTimeGraph(req, concurrentGraph)
	//cloudTests, err := provider.NewCloudTests("aws")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//for _, test := range cloudtests.Tests.TestFunctions {
	//
	//
	//	stack, err := cloudTests.GetStack(test.RequiredStack)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	err = stack.DeployStack()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	test.Fn(cloudTests.Provider, stack)
	//
	//	err =stack.RemoveStack()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
	//
	//for _, stack := range cloudTests.Stacks {
	//	err := stack.DeployStack()
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//
	//	for _, function := range stack.Functions {
	//		req, err := cloudTests.Provider.NewFunctionRequest(function.Name(), []byte(""))
	//
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		fmt.Println(function.Name())
	//		//duration := 1 * time.Second
	//		reqDelay := 20 * time.Millisecond
	//		concurrencyLimit := uint64(3)
	//		httptrace.ConcurrentRequestsSyncedOnce(req, concurrencyLimit, reqDelay)
	//	}
	//
	//	err =stack.RemoveStack()
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//
	//}
}
