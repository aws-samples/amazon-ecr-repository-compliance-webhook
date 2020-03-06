// Package function contains library units for the ecr-repository-compliance-webhook Lambda function.
// Referenced: https://github.com/kubernetes/kubernetes/blob/v1.13.0/test/images/webhook/main.go.
package function

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/ecr-repository-compliance-webhook/pkg/webhook"
)

// FunctionContainer contains the dependencies and business logic for the ecr-repository-compliance-webhook Lambda function.
type FunctionContainer struct {
	ECR ecriface.ECRAPI
}

// NewFunctionContainer creates a new FunctionContainer.
func NewFunctionContainer(ecrSvc ecriface.ECRAPI) *FunctionContainer {
	return &FunctionContainer{
		ECR: ecrSvc,
	}
}

// Handler is a type alias for the Lambda handler's function signatire.
type Handler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// GetHandler returns the function handler for ecr-repository-compliance-webhook.
func (f *FunctionContainer) GetHandler() Handler {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		request, err := webhook.NewRequestFromEvent(event)
		response := webhook.NewResponseFromRequest(request)
		if err != nil {
			log.Errorf("Error creating request from event: %v", err)
			return response.FailValidation(500, err)
		}

		pod, err := request.UnmarshalPod()
		if err != nil {
			log.Errorf("Error unmarshalling Pod: %v", err)
			return response.FailValidation(406, err)
		}
		_ = pod
		return response.PassValidation()
	}
}
