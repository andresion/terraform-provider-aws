package networkfirewall

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/networkfirewall"
)

// findLoggingConfiguration returns the LoggingConfigurationOutput from a call to DescribeLoggingConfigurationWithContext
// given the context and findFirewall ARN.
// Returns nil if the findLoggingConfiguration is not found.
func findLoggingConfiguration(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*networkfirewall.DescribeLoggingConfigurationOutput, error) {
	input := &networkfirewall.DescribeLoggingConfigurationInput{
		FirewallArn: aws.String(arn),
	}
	output, err := conn.DescribeLoggingConfigurationWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// findFirewall returns the FirewallOutput from a call to DescribeFirewallWithContext
// given the context and findFirewall ARN.
// Returns nil if the findFirewall is not found.
func findFirewall(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*networkfirewall.DescribeFirewallOutput, error) {
	input := &networkfirewall.DescribeFirewallInput{
		FirewallArn: aws.String(arn),
	}
	output, err := conn.DescribeFirewallWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// findFirewallPolicy returns the FirewallPolicyOutput from a call to DescribeFirewallPolicyWithContext
// given the context and findFirewallPolicy ARN.
// Returns nil if the findFirewallPolicy is not found.
func findFirewallPolicy(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*networkfirewall.DescribeFirewallPolicyOutput, error) {
	input := &networkfirewall.DescribeFirewallPolicyInput{
		FirewallPolicyArn: aws.String(arn),
	}
	output, err := conn.DescribeFirewallPolicyWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// findResourcePolicy returns the Policy string from a call to DescribeResourcePolicyWithContext
// given the context and resource ARN.
// Returns nil if the findResourcePolicy is not found.
func findResourcePolicy(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*string, error) {
	input := &networkfirewall.DescribeResourcePolicyInput{
		ResourceArn: aws.String(arn),
	}
	output, err := conn.DescribeResourcePolicyWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if output == nil {
		return nil, nil
	}
	return output.Policy, nil
}

// findRuleGroup returns the RuleGroupOutput from a call to DescribeRuleGroupWithContext
// given the context and findRuleGroup ARN.
// Returns nil if the findRuleGroup is not found.
func findRuleGroup(ctx context.Context, conn *networkfirewall.NetworkFirewall, arn string) (*networkfirewall.DescribeRuleGroupOutput, error) {
	input := &networkfirewall.DescribeRuleGroupInput{
		RuleGroupArn: aws.String(arn),
	}
	output, err := conn.DescribeRuleGroupWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	return output, nil
}
