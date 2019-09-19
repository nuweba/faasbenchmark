package main

import (
	"context"
	"errors"
	"strconv"
	"time"
	"os"
	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Sleep string `json:"sleep"`
}

func isWarm() bool{
	warm := os.Getenv("warm") == "true"
	os.Setenv("warm", "true")
	return warm
}

func runTest(sleepTime time.Duration) {
	time.Sleep(sleepTime)
}

func getSleepParameter(e Event) (time.Duration, error) {
	sleepTime, err := strconv.Atoi(e.Sleep)
	if err != nil || sleepTime < 0 {
		return time.Nanosecond, errors.New("invalid sleep parameter")
	}
	return time.Duration(sleepTime) * time.Millisecond, nil
}

func getParameters(e Event) (time.Duration, error) {
	return getSleepParameter(e)
}

func Handler(ctx context.Context, e Event) (map[string]interface{}, error) {
	start := time.Now()
	reused := isWarm()
	sleepTime, err := getParameters(e)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}, nil
	}
	runTest(sleepTime)
	duration := time.Since(start).Nanoseconds()
	return map[string]interface{}{"reused": reused, "duration": duration}, nil
}

func main() {
	lambda.Start(Handler)
}
