package ds

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/terraform-providers/terraform-provider-aws/internal/flex"
)

func flattenConnectSettings(
	customerDnsIps []*string,
	s *directoryservice.DirectoryConnectSettingsDescription) []map[string]interface{} {
	if s == nil {
		return nil
	}

	settings := make(map[string]interface{})

	settings["customer_dns_ips"] = flex.FlattenStringSet(customerDnsIps)
	settings["connect_ips"] = flex.FlattenStringSet(s.ConnectIps)
	settings["customer_username"] = aws.StringValue(s.CustomerUserName)
	settings["subnet_ids"] = flex.FlattenStringSet(s.SubnetIds)
	settings["vpc_id"] = aws.StringValue(s.VpcId)
	settings["availability_zones"] = flex.FlattenStringSet(s.AvailabilityZones)

	return []map[string]interface{}{settings}
}

func flattenVpcSettings(
	s *directoryservice.DirectoryVpcSettingsDescription) []map[string]interface{} {
	settings := make(map[string]interface{})

	if s == nil {
		return nil
	}

	settings["subnet_ids"] = flex.FlattenStringSet(s.SubnetIds)
	settings["vpc_id"] = aws.StringValue(s.VpcId)
	settings["availability_zones"] = flex.FlattenStringSet(s.AvailabilityZones)

	return []map[string]interface{}{settings}
}
