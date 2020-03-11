package webhook

import (
	"errors"

	v1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


var (
	// ErrMissingFailure ...
	ErrMissingFailure = errors.New("webhook: reached invalid state, no failure reaon found")
	
	// ErrBadRequest ...
	ErrBadRequest = errors.New("webhook: bad request")
	
	// BadRequestResponse ...
	BadRequestResponse = func(err error) (*v1beta1.AdmissionReview, error) {
		return &v1beta1.AdmissionReview{
			Response: &v1beta1.AdmissionResponse {
				Allowed: false,
				Result: &metav1.Status{
					Status:  metav1.StatusFailure,
					Message: err.Error(),
					Reason: metav1.StatusReasonBadRequest,
					Code:   400,
				},
			},
		}, nil
	}
)
// Response encapsulates the AdmissionResponse sent to API Gateway.
type Response struct {
	Admission *v1beta1.AdmissionResponse
}

// NewResponseFromRequest creates a Response from a Request.
// Assumes request came from Kubernetes and contains UID.
func NewResponseFromRequest(r *Request) (*Response, error){
	if r == nil || (r != nil && r.Admission == nil){
		return nil, ErrBadRequest
	}
	if r.Admission != nil && r.Admission.UID == "" {
		return nil, ErrBadRequest
	}
	return &Response{
		Admission: &v1beta1.AdmissionResponse{
			UID: r.Admission.UID,
		},
	}, nil
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
		// Need a better way to Code with Reason; maybe use gRPC code mappings?
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
		Message: "pod contains compliant ecr repositories and images",
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
