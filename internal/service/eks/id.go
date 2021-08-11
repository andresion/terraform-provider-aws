package eks

import (
	"fmt"
	"strings"
)

const addonResourceIDSeparator = ":"

func addonCreateResourceID(clusterName, addonName string) string {
	parts := []string{clusterName, addonName}
	id := strings.Join(parts, addonResourceIDSeparator)

	return id
}

func addonParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, addonResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected cluster-name%[2]saddon-name", id, addonResourceIDSeparator)
}

const fargateProfileResourceIDSeparator = ":"

func fargateProfileCreateResourceID(clusterName, fargateProfileName string) string {
	parts := []string{clusterName, fargateProfileName}
	id := strings.Join(parts, fargateProfileResourceIDSeparator)

	return id
}

func fargateProfileParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, fargateProfileResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected cluster-name%[2]sfargate-profile-name", id, fargateProfileResourceIDSeparator)
}

const identityProviderConfigResourceIDSeparator = ":"

func identityProviderConfigCreateResourceID(clusterName, configName string) string {
	parts := []string{clusterName, configName}
	id := strings.Join(parts, identityProviderConfigResourceIDSeparator)

	return id
}

func identityProviderConfigParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, identityProviderConfigResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected cluster-name%[2]sconfig-name", id, identityProviderConfigResourceIDSeparator)
}

const nodeGroupResourceIDSeparator = ":"

func nodeGroupCreateResourceID(clusterName, nodeGroupName string) string {
	parts := []string{clusterName, nodeGroupName}
	id := strings.Join(parts, nodeGroupResourceIDSeparator)

	return id
}

func nodeGroupParseResourceID(id string) (string, string, error) {
	parts := strings.Split(id, nodeGroupResourceIDSeparator)

	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%[1]s), expected cluster-name%[2]snode-group-name", id, nodeGroupResourceIDSeparator)
}
