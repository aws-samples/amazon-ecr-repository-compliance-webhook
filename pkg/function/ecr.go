package function

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
)

// RepositoryCompliant checks if the repository for the image
// that was just sent to the is compliant.
// 1. Has image tagging enabled
// 2. Has image scanning enabled
func RepositoryCompliant(ctx context.Context, repository string, svc ecriface.ECRAPI) (bool, error) {
	input := &ecr.DescribeRepositoriesInput{
		MaxResults:      aws.Int64(1),
		RepositoryNames: []*string{aws.String(repository)},
	}
	if err := input.Validate(); err != nil {
		return false, err
	}
	output, err := svc.DescribeRepositoriesWithContext(ctx, input)
	if err != nil {
		return false, err
	}
	if len(output.Repositories) == 0 {
		log.Errorf("No repositories named '%s' found:", repository)
		return false, nil
	}
	r := output.Repositories[0]
	if aws.StringValue(r.ImageTagMutability) == ecr.ImageTagMutabilityMutable || !aws.BoolValue(r.ImageScanningConfiguration.ScanOnPush) {
		log.Errorf("Error: repository '%s' is not compliant:", repository)
		return false, nil
	}
	return true, nil
}
