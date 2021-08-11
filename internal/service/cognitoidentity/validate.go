package cognitoidentity

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/service/cognitoidentity"
)

func validIdentityPoolName(v interface{}, k string) (ws []string, errors []error) {
	val := v.(string)
	if !regexp.MustCompile(`^[\w\s+=,.@-]+$`).MatchString(val) {
		errors = append(errors, fmt.Errorf("%q must contain only alphanumeric characters, dots, underscores and hyphens", k))
	}

	return
}

func validIdentityProvidersClientId(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character", k))
	}

	if len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if !regexp.MustCompile(`^[\w_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain only alphanumeric characters and underscores", k))
	}

	return
}

func validIdentityProvidersProviderName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character", k))
	}

	if len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if !regexp.MustCompile(`^[\w._:/-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain only alphanumeric characters, dots, underscores, colons, slashes and hyphens", k))
	}

	return
}

func validProviderDeveloperName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 100 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 100 characters", k))
	}

	if !regexp.MustCompile(`^[\w._-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain only alphanumeric characters, dots, underscores and hyphens", k))
	}

	return
}

func validRoleMappingsAmbiguousRoleResolutionAgainstType(v map[string]interface{}) (errors []error) {
	t := v["type"].(string)
	isRequired := t == cognitoidentity.RoleMappingTypeToken || t == cognitoidentity.RoleMappingTypeRules

	if value, ok := v["ambiguous_role_resolution"]; (!ok || value == "") && isRequired {
		errors = append(errors, fmt.Errorf(`Ambiguous Role Resolution must be defined when "type" equals "Token" or "Rules"`))
	}

	return
}

func validRoleMappingsRulesClaim(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(`^[\p{L}\p{M}\p{S}\p{N}\p{P}]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain only alphanumeric characters, dots, underscores, colons, slashes and hyphens", k))
	}

	return
}

// Validates that either authenticated or unauthenticated is defined
func validRoles(v map[string]interface{}) (errors []error) {
	k := "roles"
	_, hasAuthenticated := v["authenticated"].(string)
	_, hasUnauthenticated := v["unauthenticated"].(string)

	if !hasAuthenticated && !hasUnauthenticated {
		errors = append(errors, fmt.Errorf("%q: Either \"authenticated\" or \"unauthenticated\" must be defined", k))
	}

	return
}

func validSupportedLoginProviders(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q cannot be less than 1 character", k))
	}

	if len(value) > 128 {
		errors = append(errors, fmt.Errorf("%q cannot be longer than 128 characters", k))
	}

	if !regexp.MustCompile(`^[\w.;_/-]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("%q must contain only alphanumeric characters, dots, semicolons, underscores, slashes and hyphens", k))
	}

	return
}
