package securityhub

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
)

const (
	arnSeparator = "/"
	arnService   = "securityhub"
)

// standardsControlARNToStandardsSubscriptionARN converts a security standard control ARN to a subscription ARN.
func standardsControlARNToStandardsSubscriptionARN(inputARN string) (string, error) {
	parsedARN, err := arn.Parse(inputARN)

	if err != nil {
		return "", fmt.Errorf("error parsing ARN (%s): %w", inputARN, err)
	}

	if actual, expected := parsedARN.Service, arnService; actual != expected {
		return "", fmt.Errorf("expected service %s in ARN (%s), got: %s", expected, inputARN, actual)
	}

	inputResourceParts := strings.Split(parsedARN.Resource, arnSeparator)

	if actual, expected := len(inputResourceParts), 3; actual < expected {
		return "", fmt.Errorf("expected at least %d resource parts in ARN (%s), got: %d", expected, inputARN, actual)
	}

	outputResourceParts := append([]string{"subscription"}, inputResourceParts[1:len(inputResourceParts)-1]...)

	outputARN := arn.ARN{
		Partition: parsedARN.Partition,
		Service:   parsedARN.Service,
		Region:    parsedARN.Region,
		AccountID: parsedARN.AccountID,
		Resource:  strings.Join(outputResourceParts, arnSeparator),
	}.String()

	return outputARN, nil
}
