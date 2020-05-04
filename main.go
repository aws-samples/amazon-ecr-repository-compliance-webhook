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
	"github.com/swoldemi/amazon-ecr-repository-compliance-webhook/pkg/function"
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
