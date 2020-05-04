/*
  Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
  Licensed under the Apache License, Version 2.0 (the "License").
  You may not use this file except in compliance with the License.
  A copy of the License is located at
      http://www.apache.org/licenses/LICENSE-2.0
  or in the "license" file accompanying this file. This file is distributed
  on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
  express or implied. See the License for the specific language governing
  permissions and limitations under the License.
*/

package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	log "github.com/sirupsen/logrus"
	"github.com/aws-samples/amazon-ecr-repository-compliance-webhook/pkg/function"
)

func init() {
	xraylvl, loglvl := logLevels(os.Getenv("LOG_LEVEL"))
	log.SetFormatter(new(log.JSONFormatter))
	log.Infof("Got log levels [%s, %s]", xraylvl, loglvl)
	log.SetLevel(loglvl)
	xray.SetLogger(xraylog.NewDefaultLogger(os.Stdout, xraylvl))
}

var (
	sess = xray.AWSSession(session.Must(session.NewSession()))
	svc  = ecr.New(sess, &aws.Config{Region: getRegistryRegion()})

	// Handler is the handler for the validating webhook.
	Handler = function.NewContainer(svc).Handler().WithLogging().WithProxiedResponse()

	// Version is the shortened git hash of the binary's source code.
	// It is injected using the -X linker flag when running `make`
	Version string
)

func main() {
	log.Infof("Starting function version: %s", Version)
	lambda.Start(Handler)
}

func logLevels(lvl string) (xraylog.LogLevel, log.Level) {
	loglvl, err := log.ParseLevel(lvl)
	if err != nil {
		return xraylog.LogLevelInfo, log.InfoLevel
	}

	var xraylvl xraylog.LogLevel
	switch lvl {
	case "DEBUG":
		xraylvl = xraylog.LogLevelDebug
	case "INFO":
		xraylvl = xraylog.LogLevelInfo
	case "WARN":
		xraylvl = xraylog.LogLevelWarn
	case "ERROR":
		xraylvl = xraylog.LogLevelError
	default:
		xraylvl = xraylog.LogLevelInfo
	}
	return xraylvl, loglvl
}

func getRegistryRegion() *string {
	if value, ok := os.LookupEnv("REGISTRY_REGION"); ok {
		return aws.String(value)
	}
	return aws.String(os.Getenv("AWS_REGION"))
}
