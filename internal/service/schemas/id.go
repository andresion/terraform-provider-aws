package schemas

import (
	"fmt"
	"strings"
)

const schemaResourceIDSeparator = "/"

func schemaCreateResourceID(schemaName, registryName string) string {
	parts := []string{schemaName, registryName}
	id := strings.Join(parts, schemaResourceIDSeparator)

	return id
}

func schemaParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, schemaResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected SCHEMA_NAME%[2]sREGISTRY_NAME", id, schemaResourceIDSeparator)
}
