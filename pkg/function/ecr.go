package function

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"golang.org/x/sync/errgroup"
)

// CheckRepositoryCompliance checks if the repository for the image
// that was just sent to the is compliant.
// 1. Comes from ECR
// 2. Has image tag immutability enabled
// 3. Has image scan on push enabled
// 4. Does not contain any critical vulnerabilities
func (c *Container) CheckRepositoryCompliance(ctx context.Context, image string) (bool, error) {
	repo := strings.Split(image, ":")[0]
	input := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []*string{aws.String(repo)},
	}
	if err := input.Validate(); err != nil {
		return false, err
	}
	output, err := c.ECR.DescribeRepositoriesWithContext(ctx, input)
	if err != nil {
		return false, err
	}
	if len(output.Repositories) == 0 {
		return false, fmt.Errorf("no repositories named '%s' found", repo)
	}
	r := output.Repositories[0]
	if aws.StringValue(r.ImageTagMutability) == ecr.ImageTagMutabilityMutable {
		return false, fmt.Errorf("repository '%s' does not have image tag immutability enabled", repo)
	}
	if !aws.BoolValue(r.ImageScanningConfiguration.ScanOnPush) {
		return false, fmt.Errorf("repository '%s' does not have image scan on push enabled", repo)
	}
	critical, err := c.HasCriticalVulnerabilities(ctx, image)
	if err != nil {
		return false, err
	}
	if critical {
		return false, fmt.Errorf("image '%s' contains CRITICAL vulnerabilities", image)
	}
	return true, nil
}

// BatchCheckRepositoryCompliance checks the compliance of a given set of ECR repositories.
// False is returned if a single repository is not compliant.
func (c *Container) BatchCheckRepositoryCompliance(ctx context.Context, images []string) (bool, error) {
	g, ctx := errgroup.WithContext(ctx)
	compliances := make([]bool, len(images))

	for i, image := range images {
		i, image := i, image
		g.Go(func() error {
			compliant, err := c.CheckRepositoryCompliance(ctx, image)
			compliances[i] = compliant
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

// HasCriticalVulnerabilities checks if a container image contains for 'CRITICAL' vulnerabilities.
func (c *Container) HasCriticalVulnerabilities(ctx context.Context, image string) (bool, error) {
	var (
		segments  = strings.Split(image, ":")
		repo, tag = segments[0], segments[1]
		found     = false
	)
	input := &ecr.DescribeImageScanFindingsInput{
		ImageId: &ecr.ImageIdentifier{
			ImageTag: aws.String(tag),
		},
		RepositoryName: aws.String(repo),
	}
	if err := input.Validate(); err != nil {
		return true, err
	}

	pager := func(out *ecr.DescribeImageScanFindingsOutput, lastPage bool) bool {
		if found {
			return true // break out of paging if we've already found a critical vuln.
		}
		for _, finding := range out.ImageScanFindings.Findings {
			if aws.StringValue(finding.Severity) == ecr.FindingSeverityCritical {
				found = true
			}
		}
		return lastPage
	}

	if err := c.ECR.DescribeImageScanFindingsPagesWithContext(ctx, input, pager); err != nil {
		return true, err
	}
	return found, nil
}
