package aws

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"io/ioutil"
	"net/http"
	"net/url"
	"github.com/nuweba/httpbench/syncedtrace"
	"path"
	"time"
)

type Aws struct {
	session *session.Session
	region  string
	name    string
}

func New() (*Aws, error) {
	name := "aws"
	ses, err := session.NewSession()

	if err != nil {
		return nil, err
	}

	//todo: change
	//region, err := getRegion(ses)
	//
	//if err != nil {
	//	return nil, err
	//}
	region := "us-east-1"

	return &Aws{session: ses, region: region, name: name}, nil
}

func (aws *Aws) Name() string {
	return aws.name
}

func getRegion(session *session.Session) (string, error) {
	metaClient := ec2metadata.New(session)
	region, err := metaClient.Region()
	if err != nil {
		return "", err
	}
	return region, nil
}

func (aws *Aws) signLambdaReq(req *http.Request) error {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	seekerBody := bytes.NewReader(b)

	signer := v4.NewSigner(aws.session.Config.Credentials)
	_, err = signer.Sign(req, seekerBody, "lambda", aws.region, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (aws *Aws) buildLambdaInvokeReq(funcName string, qParams *url.Values, headers *http.Header, body *[]byte) (*http.Request, error) {
	lambdaUrl := url.URL{}

	// https://lambda.{region}.amazonaws.com/2015-03-31/functions/{functionName}/invocations?sleep={time}

	lambdaUrl.Scheme = "https"
	lambdaUrl.Host = fmt.Sprintf("lambda.%s.amazonaws.com", aws.region)
	lambdaUrl.Path = path.Join(lambdaUrl.Path, "2015-03-31/functions", funcName, "invocations")
	bodyParams := make(map[string]string)
	for qParam, val := range *qParams {
		bodyParams[qParam] = val[0]
	}
	bodyParams["body"] = base64.StdEncoding.EncodeToString(*body)

	jsonBody, err := json.Marshal(bodyParams)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", lambdaUrl.String(), ioutil.NopCloser(bytes.NewReader(jsonBody)))

	if err != nil {
		return nil, err
	}

	for k, multiH := range *headers {
		for _, h := range multiH {
			req.Header.Set(k, h)
		}
	}

	if err := aws.signLambdaReq(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (aws *Aws) NewFunctionRequest(funcName string, qParams *url.Values, headers *http.Header, body *[]byte) (func () (*http.Request, error)) {
	return func() (*http.Request, error) {
				return aws.buildLambdaInvokeReq(funcName, qParams, headers, body)
			}
}


func (aws *Aws) HttpInvocationTriggerStage() syncedtrace.TraceHookType {
	return syncedtrace.TLSHandshakeDone
}