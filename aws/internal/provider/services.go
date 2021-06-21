package provider

// Reference each service package so that each package's init runs.
import (
	_ "github.com/terraform-providers/terraform-provider-aws/aws/internal/service/amplify"
)
