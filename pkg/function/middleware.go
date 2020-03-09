package function

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/admission/v1beta1"
)

// Handler is a type alias for the Lambda handler's function signature.
type Handler func(context.Context, events.APIGatewayProxyRequest) (*v1beta1.AdmissionReview, error)

// WithLogging is a logging middleware for the Lambda handler.
func (h Handler) WithLogging() Handler {
	return func(ctx context.Context, event events.APIGatewayProxyRequest) (*v1beta1.AdmissionReview, error) {
		review, err := h(ctx, event)
		log.Infof("Responding with AdmissionReview [%+v] and error [%v]", review, err)
		return review, err
	}
}
