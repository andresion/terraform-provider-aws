package apigateway

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-aws/internal/client"
)

func buildInvokeURL(client *client.AWSClient, restApiId, stageName string) string {
	hostname := client.RegionalHostname(fmt.Sprintf("%s.execute-api", restApiId))
	return fmt.Sprintf("https://%s/%s", hostname, stageName)
}

// escapeJsonPointer escapes string per RFC 6901
// so it can be used as path in JSON patch operations
func escapeJsonPointer(path string) string {
	path = strings.Replace(path, "~", "~0", -1)
	path = strings.Replace(path, "/", "~1", -1)
	return path
}
