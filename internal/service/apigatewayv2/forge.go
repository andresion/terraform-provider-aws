package apigatewayv2

import (
	"strings"

	"github.com/terraform-providers/terraform-provider-aws/internal/hashcode"
)

// hashStringCaseInsensitive hashes strings in a case insensitive manner.
// If you want a Set of strings and are case inensitive, this is the SchemaSetFunc you want.
func hashStringCaseInsensitive(v interface{}) int {
	return hashcode.String(strings.ToLower(v.(string)))
}
