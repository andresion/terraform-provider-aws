package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func validateCognitoRoleMappingsRulesConfiguration(v map[string]interface{}) (errors []error) {
	t := v["type"].(string)
	valLength := 0
	if value, ok := v["mapping_rule"]; ok {
		valLength = len(value.([]interface{}))
	}

	if (valLength == 0) && t == cognitoidentity.RoleMappingTypeRules {
		errors = append(errors, fmt.Errorf("mapping_rule is required for Rules"))
	}

	if (valLength > 0) && t == cognitoidentity.RoleMappingTypeToken {
		errors = append(errors, fmt.Errorf("mapping_rule must not be set for Token based role mapping"))
	}

	return
}

func validExecutionFrequency() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		configservice.MaximumExecutionFrequencyOneHour,
		configservice.MaximumExecutionFrequencyThreeHours,
		configservice.MaximumExecutionFrequencySixHours,
		configservice.MaximumExecutionFrequencyTwelveHours,
		configservice.MaximumExecutionFrequencyTwentyFourHours,
	}, false)
}
