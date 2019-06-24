package azure

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench/syncedtrace"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Azure struct {
	region string
	name   string
}

func New() (*Azure, error) {
	name := "azure"

	//todo: change
	//region, err := getRegion(ses)
	//
	//if err != nil {
	//	return nil, err
	//}
	region := "East US"

	return &Azure{region: region, name: name}, nil
}

func (azure *Azure) Name() string {
	return azure.name
}

func getRegion(session *session.Session) (string, error) {
	metaClient := ec2metadata.New(session)
	region, err := metaClient.Region()
	if err != nil {
		return "", err
	}
	return region, nil
}

func (azure *Azure) buildFuncInvokeReq(funcName string, appName string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}

	// https://YOUR_REGION-YOUR_PROJECT_ID.cloudfunctions.net/FUNCTION_NAME?sleep={time}

	funcUrl.Scheme = "https"
	funcUrl.Host = fmt.Sprintf("%s.azurewebsites.net", appName)
	funcUrl.Path = fmt.Sprintf("/api/%s", funcName)

	fmt.Println(funcUrl)

	req, err := http.NewRequest("POST", funcUrl.String(), ioutil.NopCloser(bytes.NewReader(*body)))

	//req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = qParams.Encode()

	for k, multiH := range *headers {
		for _, h := range multiH {
			req.Header.Set(k, h)
		}
	}

	return req, nil
}

func (azure *Azure) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) func() (*http.Request, error) {
	return func() (*http.Request, error) {
		return azure.buildFuncInvokeReq(function.Handler(), stack.StackId(), qParams, headers, body)
	}
}

func (azure *Azure) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.TLSHandshakeDone
}
