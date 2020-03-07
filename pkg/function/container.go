// Package function contains library units for the ecr-repository-compliance-webhook Lambda function.
package function

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/swoldemi/ecr-repository-compliance-webhook/pkg/webhook"
	"k8s.io/api/admission/v1beta1"
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
type Handler func(context.Context, events.APIGatewayProxyRequest) (*v1beta1.AdmissionReview, error)

// GetHandler returns the function handler for ecr-repository-compliance-webhook.
func (c *Container) GetHandler() Handler {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (*v1beta1.AdmissionReview, error) {
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
			return logPass(response.PassValidation), nil
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
		return logPass(response.PassValidation), nil
	}
}

type failureFunc func(int32, error) (*v1beta1.AdmissionReview, error)

const code = 406

func logFailure(f failureFunc, err error) (*v1beta1.AdmissionReview, error) {
	log.Infof("Got failure code %d and error %v", code, err)
	event, err := f(code, err)
	log.Infof("Responding with failed review %#v and error %v", event, err)
	return event, err
}

func logPass(f func() *v1beta1.AdmissionReview) *v1beta1.AdmissionReview {
	event := f()
	log.Infof("Responding with passed review %#v", event)
	return event
}
