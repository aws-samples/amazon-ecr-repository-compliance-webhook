package webhook

import (
	"errors"

	v1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ErrMissingFailure ...
var ErrMissingFailure = errors.New("webhook: reached invalidate state, no failure reaon found")

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
func (r *Response) FailValidation(code int32, failure error) (*v1beta1.AdmissionReview, error) {
	if failure == nil {
		return nil, ErrMissingFailure
	}

	r.Admission.Allowed = false
	r.Admission.Result = &metav1.Status{
		Status:  metav1.StatusFailure,
		Message: failure.Error(),
		// Need a better way to Code with Reason; maybe use grpc code mappings?
		Reason: metav1.StatusReasonNotAcceptable,
		Code:   code,
	}
	return respond(r.Admission), nil
}

// PassValidation populates the AdmissionResponse with the pass contents
// (message) and returns the AdmissionReview JSON response for API Gateway.
func (r *Response) PassValidation() *v1beta1.AdmissionReview {
	r.Admission.Allowed = true
	r.Admission.Result = &metav1.Status{
		Status:  metav1.StatusSuccess,
		Message: "pod contains compliant ecr repositories",
		Code:    200,
	}
	return respond(r.Admission)
}

func respond(admission *v1beta1.AdmissionResponse) *v1beta1.AdmissionReview {
	return &v1beta1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1beta1",
		},
		Response: admission,
	}
}
