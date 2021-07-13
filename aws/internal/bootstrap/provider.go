package bootstrap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var Provider func() *schema.Provider
