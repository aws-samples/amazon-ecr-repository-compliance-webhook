// Package function contains library units for the ecr-repository-compliance-webhook Lambda function.
package function

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/ecr-repository-compliance-webhook/pkg/webhook"
)

// ErrFailedCompliance ...
var ErrFailedCompliance = errors.New("webhook: repository fails ecr criteria")

// Container contains the dependencies and business logic for the ecr-repository-compliance-webhook Lambda function.
type Container struct {
	ECR ecriface.ECRAPI
}

// NewContainer creates a new function Container.
func NewContainer(ecrSvc ecriface.ECRAPI) *Container {
	return &Container{
		ECR: ecrSvc,
	}
}

// Handler is a type alias for the Lambda handler's function signatire.
type Handler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// GetHandler returns the function handler for ecr-repository-compliance-webhook.
func (c *Container) GetHandler() Handler {
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

		if webhook.InCriticalNamespace(pod) {
			log.Info("Pod is in critical namespace, automatically passing")
			return response.PassValidation()
		}

		repos, err := webhook.ParseRepositories(pod)
		if err != nil {
			log.Errorf("Error extracting repositories: %v", err)
			return response.FailValidation(500, err)
		}

		compliant, err := c.BatchCheckRepositoryCompliance(ctx, repos)
		if err != nil {
			log.Errorf("Error during compliance check: %v", err)
			return response.FailValidation(500, err)
		}

		if !compliant {
			return response.FailValidation(403, ErrFailedCompliance)
		}
		return response.PassValidation()
	}
}
