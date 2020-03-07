package function

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"golang.org/x/sync/errgroup"
)

// CheckRepositoryCompliance checks if the repository for the image
// that was just sent to the is compliant.
// 1. Comes from ECR
// 2. Has image tagging enabled
// 3. Has image scanning enabled
func (c *Container) CheckRepositoryCompliance(ctx context.Context, repository string) (bool, error) {
	input := &ecr.DescribeRepositoriesInput{
		MaxResults:      aws.Int64(1),
		RepositoryNames: []*string{aws.String(repository)},
	}
	if err := input.Validate(); err != nil {
		return false, err
	}
	output, err := c.ECR.DescribeRepositoriesWithContext(ctx, input)
	if err != nil {
		return false, err
	}
	if len(output.Repositories) == 0 {
		log.Errorf("No repositories named '%s' found", repository)
		return false, nil
	}
	r := output.Repositories[0]
	if aws.StringValue(r.ImageTagMutability) == ecr.ImageTagMutabilityMutable || !aws.BoolValue(r.ImageScanningConfiguration.ScanOnPush) {
		log.Errorf("Repository '%s' is not compliant", repository)
		return false, nil
	}
	return true, nil
}

// BatchCheckRepositoryCompliance checks the compliance of a given set of ECR repositories.
// False is returned if a single repository is not compliant.
func (c *Container) BatchCheckRepositoryCompliance(ctx context.Context, repos []string) (bool, error) {
	g, ctx := errgroup.WithContext(ctx)
	compliances := make([]bool, len(repos))

	for i, repo := range repos {
		i, repo := i, repo
		g.Go(func() error {
			compliant, err := c.CheckRepositoryCompliance(ctx, repo)
			if err == nil {
				compliances[i] = compliant
			}
			return err
		})
	}
	if err := g.Wait(); err != nil {
		return false, err
	}

	for _, compliance := range compliances {
		if !compliance {
			return false, nil
		}
	}
	return true, nil
}
