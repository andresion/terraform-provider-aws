package route53

import (
	"fmt"
	"strings"
)

const keySigningKeyResourceIDSeparator = ","

func keySigningKeyCreateResourceID(transitGatewayRouteTableID string, prefixListID string) string {
	parts := []string{transitGatewayRouteTableID, prefixListID}
	id := strings.Join(parts, keySigningKeyResourceIDSeparator)

	return id
}

func keySigningKeyParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, keySigningKeyResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected hosted-zone-id%[2]sname", id, keySigningKeyResourceIDSeparator)
}
