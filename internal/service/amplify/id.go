package amplify

import (
	"fmt"
	"strings"
)

const backendEnvironmentResourceIDSeparator = "/"

func backendEnvironmentCreateResourceID(appID, environmentName string) string {
	parts := []string{appID, environmentName}
	id := strings.Join(parts, backendEnvironmentResourceIDSeparator)

	return id
}

func backendEnvironmentParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, backendEnvironmentResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected APPID%[2]sENVIRONMENTNAME", id, backendEnvironmentResourceIDSeparator)
}

const branchResourceIDSeparator = "/"

func branchCreateResourceID(appID, branchName string) string {
	parts := []string{appID, branchName}
	id := strings.Join(parts, branchResourceIDSeparator)

	return id
}

func branchParseResourceID(id string) (string, string, error) {
	parts := strings.SplitN(id, branchResourceIDSeparator, 2)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected APPID%[2]sBRANCHNAME", id, branchResourceIDSeparator)
}

const domainAssociationResourceIDSeparator = "/"

func domainAssociationCreateResourceID(appID, domainName string) string {
	parts := []string{appID, domainName}
	id := strings.Join(parts, domainAssociationResourceIDSeparator)

	return id
}

func domainAssociationParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, domainAssociationResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected APPID%[2]sDOMAINNAME", id, domainAssociationResourceIDSeparator)
}
