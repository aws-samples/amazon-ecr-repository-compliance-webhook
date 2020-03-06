package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type mockECRClient struct {
	mock.Mock
	ecriface.ECRAPI
}

// CreateSnapshot mocks the PutMetricData CloudWatch API endpoint.
func (_m *mockECRClient) DescribeRepositoriesWithContext(ctx aws.Context, input *ecr.DescribeRepositoriesInput, opts ...request.Option) (*ecr.DescribeRepositoriesOutput, error) {
	log.Debugf("Mocking PutMetricData API with input: %s\n", input.String())
	args := _m.Called(ctx, input)
	return args.Get(0).(*ecr.DescribeRepositoriesOutput), args.Error(1)
}

func TestHandler(t *testing.T) {
	ecrSvc := new(mockECRClient)
	ecrSvc.On("DescribeRepositoriesWithContext", context.Background(), &ecr.DescribeRepositoriesInput{}).
		Return(&ecr.DescribeRepositoriesOutput{}, nil)
	ecrSvc.AssertExpectations(t)
}
