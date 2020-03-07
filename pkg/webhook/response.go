package webhook

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	v1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	defaultHeaders = map[string]string{"Content-Type": "application/json"}

	// ErrMissingFailure ...
	ErrMissingFailure = errors.New("webhook: reached invalidate state, no failure reaon found")
)

// Response encapsulates the AdmissionResponse sent to API Gateway
type Response struct {
	Admission *v1beta1.AdmissionResponse
}

// NewResponseFromRequest creates a Response from a Request.
func NewResponseFromRequest(r *Request) *Response {
	// TODO: What's the best way to handle no UID?
	return &Response{
		Admission: &v1beta1.AdmissionResponse{
			UID: r.Admission.UID,
		},
	}
}

// FailValidation populates the AdmissionResponse with the failure contents
// (message and error) and returns the AdmissionReview JSON body response for API Gateway.
func (r *Response) FailValidation(code int32, failure error) (events.APIGatewayProxyResponse, error) {
	out := events.APIGatewayProxyResponse{
		Headers:    defaultHeaders,
		StatusCode: int(code),
	}
	if failure == nil {
		out.Body = ErrMissingFailure.Error()
		return out, ErrMissingFailure
	}

	r.Admission.Allowed = false
	r.Admission.Result = &metav1.Status{
		Status:  metav1.StatusFailure,
		Message: failure.Error(),
		// Need a better way to Code with Reason; maybe use grpc code mappings?
		Reason: metav1.StatusReasonUnknown,
		Code:   code,
	}
	body, err := marshalResponse(r.Admission)
	if err != nil {
		out.Body = err.Error()
		return out, err
	}
	out.Body = string(body)
	return out, nil
}

// PassValidation populates the AdmissionResponse with the pass contents
// (message) and returns the AdmissionReview JSON response for API Gateway.
func (r *Response) PassValidation() (events.APIGatewayProxyResponse, error) {
	r.Admission.Allowed = true
	r.Admission.Result = &metav1.Status{
		Status:  metav1.StatusSuccess,
		Message: "pod contains compliant ecr repositories",
		Code:    200,
	}
	out := events.APIGatewayProxyResponse{
		Headers: defaultHeaders,
	}

	body, err := marshalResponse(r.Admission)
	if err != nil {
		out.StatusCode = 500
		out.Body = err.Error()
		return out, err
	}
	out.StatusCode = 200
	out.Body = string(body)
	return out, nil
}

func marshalResponse(admission *v1beta1.AdmissionResponse) ([]byte, error) {
	review := &v1beta1.AdmissionReview{
		Response: admission,
	}
	out, err := json.Marshal(review)
	if err != nil {
		return nil, err
	}
	return out, nil
}
