package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-xray-sdk-go/xray"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/amazon-ecr-repository-compliance-webhook/pkg/function"
)

func getRegistryRegion() *string {
	if value, ok := os.LookupEnv("REGISTRY_REGION"); ok {
		return aws.String(value)
	}
	return aws.String(os.Getenv("AWS_REGION"))
}

var (
	sess = xray.AWSSession(session.Must(session.NewSession()))
	svc  = ecr.New(sess, &aws.Config{Region: getRegistryRegion()})

	// Handler is the handler for the validating webhook.
	Handler = function.NewContainer(svc).Handler().WithLogging().WithProxiedResponse()

	// Version is the shorted git hash of the binary's source code.
	// It is injected using the -X linker flag when running `make`
	Version string
)

func main() {
	log.Infof("Starting function version: %s", Version)
	lambda.Start(Handler)
}
