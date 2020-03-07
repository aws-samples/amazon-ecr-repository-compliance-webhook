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

var (
	// ErrFailedCompliance ...
	ErrFailedCompliance = errors.New("webhook: repository fails ecr criteria")

	// ErrImagesNotFound ...
	ErrImagesNotFound = errors.New("webhook: no ecr images found in pod specification")
)

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
			return logFailure(response.FailValidation, err)
		}

		pod, err := request.UnmarshalPod()
		if err != nil {
			log.Errorf("Error unmarshalling Pod: %v", err)
			return logFailure(response.FailValidation, err)
		}

		if webhook.InCriticalNamespace(pod) {
			log.Info("Pod is in critical namespace, automatically passing")
			return logPass(response.PassValidation)
		}

		repos, err := webhook.ParseRepositories(pod)
		if err != nil {
			log.Errorf("Error extracting repositories: %v", err)
			return logFailure(response.FailValidation, err)
		}
		if len(repos) == 0 {
			return logFailure(response.FailValidation, ErrImagesNotFound)
		}

		compliant, err := c.BatchCheckRepositoryCompliance(ctx, repos)
		if err != nil {
			log.Errorf("Error during compliance check: %v", err)
			return logFailure(response.FailValidation, err)
		}

		if !compliant {
			return logFailure(response.FailValidation, ErrFailedCompliance)
		}
		return logPass(response.PassValidation)
	}
}

type failureFunc func(int32, error) (events.APIGatewayProxyResponse, error)

const code = 406

func logFailure(f failureFunc, err error) (events.APIGatewayProxyResponse, error) {
	log.Infof("Got failure code %d and error %v", code, err)
	event, err := f(code, err)
	log.Infof("Responding with failure event %#v and error %v", event, err)
	return event, err
}

func logPass(f func() (events.APIGatewayProxyResponse, error)) (events.APIGatewayProxyResponse, error) {
	event, err := f()
	log.Infof("Responding with pass event %#v and error %v", event, err)
	return event, err
}
