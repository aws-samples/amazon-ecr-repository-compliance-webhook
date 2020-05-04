// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/swoldemi/amazon-ecr-repository-compliance-webhook/pkg/function"
	"github.com/swoldemi/amazon-ecr-repository-compliance-webhook/testdata"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mockECRClient struct {
	mock.Mock
	ecriface.ECRAPI
}

// DescribeRepositoriesWithContext mocks the DescribeRepositories ECR API endpoint.
func (_m *mockECRClient) DescribeRepositoriesWithContext(ctx aws.Context, input *ecr.DescribeRepositoriesInput, opts ...request.Option) (*ecr.DescribeRepositoriesOutput, error) {
	log.Debugf("Mocking DescribeRepositories API with input: %s\n", input.String())
	args := _m.Called(ctx, input)
	return args.Get(0).(*ecr.DescribeRepositoriesOutput), args.Error(1)
}

// DescribeImageScanFindingsPagesWithContext mocks the DescribeImageScanFindingsP ECR API endpoint.
func (_m *mockECRClient) DescribeImageScanFindingsPagesWithContext(ctx aws.Context, input *ecr.DescribeImageScanFindingsInput, fn func(*ecr.DescribeImageScanFindingsOutput, bool) bool, opts ...request.Option) error {
	log.Debugf("Mocking DescribeImageScanFindings API with input: %s\n", input.String())
	args := _m.Called(ctx, input, fn)
	return args.Error(0)
}

func TestHandler(t *testing.T) {
	type args struct {
		image           string
		repo            *ecr.Repository
		shouldCheckVuln bool
		scanFindings    *ecr.DescribeImageScanFindingsOutput
		event           events.APIGatewayProxyRequest
	}
	ecrSvc := new(mockECRClient)
	h := function.NewContainer(ecrSvc).Handler()

	tests := []struct {
		name    string
		args    args
		status  string
		wantErr bool
	}{
		{
			name: "BadRequestFailure",
			args: args{
				shouldCheckVuln: false,
				repo:            nil,
				event:           eventWithBadRequest(),
			},
			status:  metav1.StatusFailure,
			wantErr: true,
		},
		{
			name: "BadRequestNoUIDFailure",
			args: args{
				shouldCheckVuln: false,
				repo:            nil,
				event:           eventWithNoUID(),
			},
			status:  metav1.StatusFailure,
			wantErr: true,
		},
		{
			name: "ImmutabilityAndScanningDisabledFailure",
			args: args{
				image:           "auth:notlatest",
				shouldCheckVuln: false,
				repo: &ecr.Repository{
					RepositoryName:             aws.String("auth"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityMutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(false)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/auth:notlatest"),
			},
			status:  metav1.StatusFailure,
			wantErr: true,
		},
		{
			name: "ScanningDisabledFailure",
			args: args{
				image:           "shipping:notlatest",
				shouldCheckVuln: false,
				repo: &ecr.Repository{
					RepositoryName:             aws.String("shipping"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityImmutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(false)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/shipping:notlatest"),
			},
			status:  metav1.StatusFailure,
			wantErr: true,
		},
		{
			name: "ImmutabilityDisabledFailure",
			args: args{
				image:           "ordering:notlatest",
				shouldCheckVuln: false,
				repo: &ecr.Repository{
					RepositoryName:             aws.String("ordering"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityMutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(true)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/ordering:notlatest"),
			},
			status:  metav1.StatusFailure,
			wantErr: true,
		},
		{
			name: "NotECRRepositoryFailure",
			args: args{
				image:           "nginx-ingress-controller:0.30.0",
				shouldCheckVuln: false,
				repo:            nil,
				event:           eventWithImage("quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.30.0"),
			},
			status:  metav1.StatusFailure,
			wantErr: true,
		},
		{
			name: "HasCriticalVulnerabilitiesFailure",
			args: args{
				image:           "fakecompany/exploitedpush:notlatest",
				shouldCheckVuln: true,
				scanFindings:    findingsWithCriticalVuln(),
				repo: &ecr.Repository{
					RepositoryName:             aws.String("fakecompany/exploitedpush"),
					ImageTagMutability:         aws.String(ecr.ImageTagMutabilityImmutable),
					ImageScanningConfiguration: &ecr.ImageScanningConfiguration{ScanOnPush: aws.Bool(true)},
				},
				event: eventWithImage("123456789012.dkr.ecr.region.amazonaws.com/fakecompany/exploitedpush:notlatest"),
			},
			status:  metav1.StatusFailure,
			wantErr: true,
		},
		{
			name: "ImmutabilityAndScanningEnabledAndNoCriticalVulnerabilitiesPass",
			args: args{
				image:           "costcenter1/payroll:notlatest",
				shouldCheckVuln: true,
				scanFindings:    findingsWithNoVuln(),
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
						RepositoryNames: []*string{tt.args.repo.RepositoryName},
					},
				).Return(&ecr.DescribeRepositoriesOutput{Repositories: []*ecr.Repository{tt.args.repo}}, nil)
			}

			if tt.args.shouldCheckVuln {
				// DescribeImageScanFindingsPagesWithContext will not be called if
				// the repository fails one of the first three checks
				tag := strings.Split(tt.args.image, ":")[1]
				ecrSvc.On("DescribeImageScanFindingsPagesWithContext",
					ctx,
					&ecr.DescribeImageScanFindingsInput{
						ImageId: &ecr.ImageIdentifier{
							ImageTag: aws.String(tag),
						},
						RepositoryName: tt.args.repo.RepositoryName,
					},
					mock.AnythingOfType("func(*ecr.DescribeImageScanFindingsOutput, bool) bool"),
				).Return(nil).Run(func(args mock.Arguments) {
					arg := args.Get(2).(func(*ecr.DescribeImageScanFindingsOutput, bool) bool)
					arg(tt.args.scanFindings, true)
				})
			}
			review, err := h(context.Background(), tt.args.event)
			if err != nil {
				t.Fatalf("Error during request for image: %v", err)
			}
			t.Logf("Got review body: %#+v", review)
			require.Nil(t, err)
			require.Equal(t, tt.status, review.Response.Result.Status)
			if tt.wantErr {
				require.GreaterOrEqual(t, review.Response.Result.Code, int32(400))
				require.Less(t, review.Response.Result.Code, int32(500))
			}
			ecrSvc.AssertExpectations(t)
		})
	}
}

var base = events.APIGatewayProxyRequest{
	Path:       "/check-image-compliance",
	HTTPMethod: "POST",
	Headers:    map[string]string{"content-type": "application/json"},
}

func eventWithNoUID() events.APIGatewayProxyRequest {
	base.Body = testdata.ReviewWithNoUID
	return base
}

func eventWithBadRequest() events.APIGatewayProxyRequest {
	base.Body = testdata.ReviewWithBadRequest
	return base
}

func eventWithImage(image string) events.APIGatewayProxyRequest {
	base.Body = fmt.Sprintf(testdata.ReviewWithOneImage, image)
	return base
}

func findingsWithCriticalVuln() *ecr.DescribeImageScanFindingsOutput {
	return &ecr.DescribeImageScanFindingsOutput{
		ImageScanFindings: &ecr.ImageScanFindings{
			Findings: []*ecr.ImageScanFinding{
				{
					Severity: aws.String(ecr.FindingSeverityCritical),
				},
			},
		},
	}
}

func findingsWithNoVuln() *ecr.DescribeImageScanFindingsOutput {
	return &ecr.DescribeImageScanFindingsOutput{
		ImageScanFindings: &ecr.ImageScanFindings{
			Findings: []*ecr.ImageScanFinding{
				{
					Severity: aws.String(ecr.FindingSeverityInformational),
				},
			},
		},
	}
}
