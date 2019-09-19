package google

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/gddo/httputil/header"
	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/httpbench/syncedtrace"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

type Google struct {
	region string
	name   string
}

func New() (*Google, error) {
	name := "google"

	//todo: change
	//region, err := getRegion(ses)
	//
	//if err != nil {
	//	return nil, err
	//}
	region := "us-central1"

	return &Google{region: region, name: name}, nil
}

func (google *Google) Name() string {
	return google.name
}

func getRegion(session *session.Session) (string, error) {
	metaClient := ec2metadata.New(session)
	region, err := metaClient.Region()
	if err != nil {
		return "", err
	}
	return region, nil
}

func (google *Google) buildGFuncInvokeReq(funcName string, projectId string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	funcUrl := url.URL{}

	// https://YOUR_REGION-YOUR_PROJECT_ID.cloudfunctions.net/FUNCTION_NAME?sleep={time}

	funcUrl.Scheme = "https"
	funcUrl.Host = fmt.Sprintf("%s-%s.cloudfunctions.net", google.region, projectId)
	funcUrl.Path = path.Join(funcUrl.Path, funcName)

	req, err := http.NewRequest("POST", funcUrl.String(), ioutil.NopCloser(bytes.NewReader(*body)))

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

func (google *Google) NewFunctionRequest(stack stack.Stack, function stack.Function, qParams *url.Values, headers *http.Header, body *[]byte) (func(uniqueId string) (*http.Request, error)) {
	return func(uniqueId string) (*http.Request, error) {
		localHeaders := header.Copy(*headers)
		localHeaders.Add("Faastest-id", uniqueId)
		return google.buildGFuncInvokeReq(function.Handler(),stack.Project(), qParams, &localHeaders, body)
	}
}

func (google *Google) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.TLSHandshakeDone
}
