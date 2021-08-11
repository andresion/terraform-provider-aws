package cloudwatchevents

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	partnerEventBusPattern = regexp.MustCompile(`^aws\.partner(/[\.\-_A-Za-z0-9]+){2,}$`)
)

const defaultEventBusName = "default"

const permissionIDSeparator = "/"

func permissionCreateID(eventBusName, statementID string) string {
	if eventBusName == "" || eventBusName == defaultEventBusName {
		return statementID
	}
	return eventBusName + permissionIDSeparator + statementID
}

func permissionParseID(id string) (string, string, error) {
	parts := strings.Split(id, permissionIDSeparator)
	if len(parts) == 1 && parts[0] != "" {
		return defaultEventBusName, parts[0], nil
	}
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected format for ID (%q), expected <event-bus-name>"+permissionIDSeparator+"<statement-id> or <statement-id>", id)
}

const ruleIDSeparator = "/"

func ruleCreateID(eventBusName, ruleName string) string {
	if eventBusName == "" || eventBusName == defaultEventBusName {
		return ruleName
	}
	return eventBusName + ruleIDSeparator + ruleName
}

func ruleParseID(id string) (string, string, error) {
	parts := strings.Split(id, ruleIDSeparator)
	if len(parts) == 1 && parts[0] != "" {
		return defaultEventBusName, parts[0], nil
	}
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return parts[0], parts[1], nil
	}
	if len(parts) > 2 {
		i := strings.LastIndex(id, ruleIDSeparator)
		busName := id[:i]
		statementID := id[i+1:]
		if partnerEventBusPattern.MatchString(busName) && statementID != "" {
			return busName, statementID, nil
		}
	}

	return "", "", fmt.Errorf("unexpected format for ID (%q), expected <event-bus-name>"+ruleIDSeparator+"<rule-name> or <rule-name>", id)
}

// Terraform state IDs for Targets are not parseable, since the separator used ("-") is also a valid
// character in both the rule name and the target id.

const targetIDSeparator = "-"
const targetImportIDSeparator = "/"

func targetCreateID(eventBusName, ruleName, targetID string) string {
	id := ruleName + targetIDSeparator + targetID
	if eventBusName != "" && eventBusName != defaultEventBusName {
		id = eventBusName + targetIDSeparator + id
	}
	return id
}

func targetParseImportID(id string) (string, string, string, error) {
	parts := strings.Split(id, targetImportIDSeparator)
	if len(parts) == 2 && parts[0] != "" && parts[1] != "" {
		return defaultEventBusName, parts[0], parts[1], nil
	}
	if len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != "" {
		return parts[0], parts[1], parts[2], nil
	}
	if len(parts) > 3 {
		iTarget := strings.LastIndex(id, targetImportIDSeparator)
		targetID := id[iTarget+1:]
		iRule := strings.LastIndex(id[:iTarget], targetImportIDSeparator)
		busName := id[:iRule]
		ruleName := id[iRule+1 : iTarget]
		if partnerEventBusPattern.MatchString(busName) && ruleName != "" && targetID != "" {
			return busName, ruleName, targetID, nil
		}
	}

	return "", "", "", fmt.Errorf("unexpected format for ID (%q), expected <event-bus-name>"+targetImportIDSeparator+"<rule-name>"+targetImportIDSeparator+"<target-id> or <rule-name>"+targetImportIDSeparator+"<target-id>", id)
}
