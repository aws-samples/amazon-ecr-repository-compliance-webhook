package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/swoldemi/ecr-repository-compliance-webhook/pkg/function"
	"github.com/swoldemi/ecr-repository-compliance-webhook/pkg/webhook"
	"github.com/swoldemi/ecr-repository-compliance-webhook/testdata"
	v1beta1 "k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mockECRClient struct {
	mock.Mock
	ecriface.ECRAPI
}

// DescribeRepositoriesWithContext mocks the DescribeRepositories ECR  API endpoint.
func (_m *mockECRClient) DescribeRepositoriesWithContext(ctx aws.Context, input *ecr.DescribeRepositoriesInput, opts ...request.Option) (*ecr.DescribeRepositoriesOutput, error) {
	log.Debugf("Mocking PutMetricData API with input: %s\n", input.String())
	args := _m.Called(ctx, input)
	return args.Get(0).(*ecr.DescribeRepositoriesOutput), args.Error(1)
}

func TestHandler(t *testing.T) {
	type args struct {
		image string
		repo  *ecr.Repository
		event events.APIGatewayProxyRequest
	}
	ecrSvc := new(mockECRClient)
	h := function.NewContainer(ecrSvc).GetHandler()

	tests := []struct {
		name    string
		args    args
		status  string
		wantErr bool
	}{
		{
			name: "ImmutabilityAndScanningDisabledFailure",
			args: args{
				image: "auth",
				repo: &ecr.Repository{
					RepositoryName:             aws.String("auth"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityMutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(false)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/auth:notlatest"),
			},
			status:  metav1.StatusFailure,
			wantErr: false,
		},
		{
			name: "ScanningDisabledFailure",
			args: args{
				image: "shipping",
				repo: &ecr.Repository{
					RepositoryName:             aws.String("shipping"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityImmutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(false)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/shipping:notlatest"),
			},
			status:  metav1.StatusFailure,
			wantErr: false,
		},
		{
			name: "ImmutabilityDisabledFailure",
			args: args{
				image: "shipping",
				repo: &ecr.Repository{
					RepositoryName:             aws.String("shipping"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityMutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(true)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/shipping:notlatest"),
			},
			status:  metav1.StatusFailure,
			wantErr: false,
		},
		{
			name: "NotECRRepositoryFailure",
			args: args{
				image: "nginx-ingress-controller",
				repo:  nil,
				event: eventWithImage("quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.30.0"),
			},
			status:  metav1.StatusFailure,
			wantErr: false,
		},
		{
			name: "ImmutabilityAndScanningEnabledPass",
			args: args{
				image: "costcenter1/payroll",
				repo: &ecr.Repository{
					RepositoryName:             aws.String("costcenter1/payroll"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityImmutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(true)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/costcenter1/payroll:notlatest"),
			},
			status:  metav1.StatusSuccess,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			if tt.args.repo != nil {
				ecrSvc.On("DescribeRepositoriesWithContext",
					ctx,
					&ecr.DescribeRepositoriesInput{
						MaxResults:      aws.Int64(1),
						RepositoryNames: []*string{tt.args.repo.RepositoryName},
					},
				).Return(&ecr.DescribeRepositoriesOutput{Repositories: []*ecr.Repository{tt.args.repo}}, nil)
			}

			response, err := h(context.Background(), tt.args.event)
			if err != nil {
				t.Fatalf("Error during request for image: %v", err)
			}
			t.Logf("Got response body: %#+v", response.Body)

			require.Nil(t, err)
			if tt.status == metav1.StatusFailure {
				require.NotEqual(t, response.StatusCode, 200)
			}

			var admission v1beta1.AdmissionReview
			if err := webhook.DeserializeReview(response.Body, &admission); err != nil {
				t.Fatalf("Error deseralizing AdmissionReview: %v", err)
			}
			require.Equal(t, tt.status, admission.Response.Result.Status)
			ecrSvc.AssertExpectations(t)
		})
	}
}

func eventWithImage(image string) events.APIGatewayProxyRequest {
	return events.APIGatewayProxyRequest{
		Path:       "/check-image-compliance",
		HTTPMethod: "POST",
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       fmt.Sprintf(testdata.ReviewWithOneImage, image),
	}
}
