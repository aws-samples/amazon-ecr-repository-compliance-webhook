// Package function contains library units for the amazon-ecr-repository-compliance-webhook Lambda function.
package function

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/amazon-ecr-repository-compliance-webhook/pkg/webhook"
	"k8s.io/api/admission/v1beta1"
)

var (
	// ErrFailedCompliance ...
	ErrFailedCompliance = errors.New("webhook: repository fails ecr criteria")

	// ErrImagesNotFound ...
	ErrImagesNotFound = errors.New("webhook: no ecr images found in pod specification")
)

// Container contains the dependencies and business logic for the amazon-ecr-repository-compliance-webhook Lambda function.
type Container struct {
	ECR ecriface.ECRAPI
}

// NewContainer creates a new function Container.
func NewContainer(ecrSvc ecriface.ECRAPI) *Container {
	return &Container{
		ECR: ecrSvc,
	}
}

// default HTTP status code to return on rejected admission
const code = 406 // NotAcceptable

// GetHandler returns the function handler for the amazon-ecr-repository-compliance-webhook.
func (c *Container) GetHandler() Handler {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (*v1beta1.AdmissionReview, error) {
		request, err := webhook.NewRequestFromEvent(event)
		if err != nil {
			log.Errorf("Error creating request from event: %v", err)
			return webhook.BadRequestResponse(err)
		}

		response, err := webhook.NewResponseFromRequest(request)
		if err != nil {
			log.Errorf("Error crafting response from request: %v", err)
			return webhook.BadRequestResponse(err)
		}

		pod, err := request.UnmarshalPod()
		if err != nil {
			log.Errorf("Error unmarshalling Pod: %v", err)
			return response.FailValidation(code, err)
		}

		if webhook.InCriticalNamespace(pod) {
			log.Info("Pod is in critical namespace, automatically passing")
			return response.PassValidation(), nil
		}

		repos, err := webhook.ParseRepositories(pod)
		if err != nil {
			log.Errorf("Error extracting repositories: %v", err)
			return response.FailValidation(code, err)
		}
		if len(repos) == 0 {
			return response.FailValidation(code, ErrImagesNotFound)
		}

		compliant, err := c.BatchCheckRepositoryCompliance(ctx, repos)
		if err != nil {
			log.Errorf("Error during compliance check: %v", err)
			return response.FailValidation(code, err)
		}

		if !compliant {
			return response.FailValidation(code, ErrFailedCompliance)
		}
		return response.PassValidation(), nil
	}
}
