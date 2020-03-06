package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-xray-sdk-go/xray"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/ecr-repository-compliance-webhook/pkg/function"
)

func main() {
	log.Info("Starting Lambda in live environment")
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
			SharedConfigState: session.SharedConfigEnable,
		},
	)
	if err != nil {
		log.Fatalf("Error creating session: %v\n", err)
		return
	}

	ecrSvc := ecr.New(sess)
	if err := xray.Configure(xray.Config{LogLevel: "trace"}); err != nil {
		log.Fatalf("Error configuring X-Ray: %v\n", err)
		return
	}

	xray.AWS(ecrSvc.Client)
	if err != nil {
		log.Fatalf("Error initializing Lambda environment: %v\n", err)
		return
	}
	lambda.Start(function.NewFunctionContainer(ecrSvc).GetHandler())
}
