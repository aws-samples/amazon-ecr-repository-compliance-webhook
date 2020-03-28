package webhook

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	v1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

const amazonawscom = "amazonaws.com"

// Errors returned when a request or resource expectation fails.
var (
	ErrInvalidContentType = errors.New("webhook: invalid content type; expected application/json")
	ErrMissingContentType = errors.New("webhook: missing content-type header")
	ErrObjectNotFound     = errors.New("webhook: request did not include object")
	ErrUnexpectedResource = errors.New("webhook: expected pod resource")
	ErrInvalidAdmission   = errors.New("webhook: admission request was nil")
	ErrInvalidContainer   = errors.New("webhook: container image is not from ecr")
)

var (
	runtimeScheme     = runtime.NewScheme()
	codecs            = serializer.NewCodecFactory(runtimeScheme)
	deserializer      = codecs.UniversalDeserializer()
	ignoredNamespaces = []string{
		metav1.NamespaceSystem,
	}
	podDefault = &schema.GroupVersionKind{
		Group:   "core",
		Version: "v1",
		Kind:    "Pod",
	}
)

// Request encapsulates the AdmissionRequest from the
// AdmissionReview proxied to the Lambda function.
type Request struct {
	Admission *v1beta1.AdmissionRequest
}

// NewRequestFromEvent creates a Request from the APIGatewayProxyRequest.
func NewRequestFromEvent(event events.APIGatewayProxyRequest) (*Request, error) {
	val, ok := event.Headers["content-type"]
	if !ok {
		return nil, ErrMissingContentType
	}
	if val != "application/json" {
		return nil, ErrInvalidContentType
	}
	var review v1beta1.AdmissionReview
	if _, _, err := deserializer.Decode([]byte(event.Body), nil, &review); err != nil {
		return nil, err
	}
	return &Request{Admission: review.Request}, nil
}

// UnmarshalPod unmarshals the raw object in the AdmissionRequest into a Pod.
func (r *Request) UnmarshalPod() (*corev1.Pod, error) {
	if r.Admission == nil {
		return nil, ErrInvalidAdmission
	}
	if len(r.Admission.Object.Raw) == 0 {
		return nil, ErrObjectNotFound
	}
	if r.Admission.Kind.Kind != podDefault.Kind {
		// If the ValidatingWebhookConfiguration was given additional resource scopes.
		return nil, ErrUnexpectedResource
	}

	var pod corev1.Pod
	if err := json.Unmarshal(r.Admission.Object.Raw, &pod); err != nil {
		return nil, err
	}
	return &pod, nil
}

// InCriticalNamespace checks that the request was for a resource
// that is being deployed into a cricial name space; e.g. kube-system.
func InCriticalNamespace(pod *corev1.Pod) bool {
	for _, n := range ignoredNamespaces {
		if pod.Namespace == n {
			return true
		}
	}
	return false
}

// ParseImages returns the container images in the Pod spec
// that originate from an Amazon ECR repository.
func ParseImages(pod *corev1.Pod) []string {
	var (
		images     []string
		containers = append(pod.Spec.Containers, pod.Spec.InitContainers...)
	)
	for _, c := range containers {
		parsed := parse(c.Image)
		if c.Image != "" && strings.Contains(c.Image, amazonawscom) && parsed != "" && !contains(images, parsed) {
			images = append(images, parsed)
		}
	}
	return images
}

// From aws_account_id.dkr.ecr.aws_region.amazonaws.com/repository:tag to repository:tag
// Or aws_account_id.dkr.ecr.aws_region.amazonaws.com/repository@sha256:hash to repository@sha256:hash
func parse(image string) string {
	if !strings.Contains(image, "/") {
		return ""
	}
	base := strings.SplitN(image, "/", 2)
	if len(base) < 1 {
		return ""
	}
	return base[1]
}

// check that we already haven't seen this image before
func contains(images []string, image string) bool {
	for _, r := range images {
		if r == image {
			return true
		}
	}
	return false
}
