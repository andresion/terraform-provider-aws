package ec2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tfnet "github.com/terraform-providers/terraform-provider-aws/internal/net"
)

// suppressEqualCIDRBlockDiffs provides custom difference suppression for CIDR blocks
// that have different string values but represent the same CIDR.
func suppressEqualCIDRBlockDiffs(k, old, new string, d *schema.ResourceData) bool {
	return tfnet.CIDRBlocksEqual(old, new)
}
