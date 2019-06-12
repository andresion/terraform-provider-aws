package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Many of the Lex resources require complex nested objects. Terraform maps only support simple key
// value pairs and not complex or mixed types. That is why these resources are defined using the
// schema.TypeList and a max of 1 item instead of the schema.TypeMap.

func expandLexSet(s *schema.Set) (items []map[string]interface{}) {
	for _, rawItem := range s.List() {
		item, ok := rawItem.(map[string]interface{})
		if !ok {
			continue
		}

		items = append(items, item)
	}

	return
}

func flattenLexEnumerationValues(values []*lexmodelbuildingservice.EnumerationValue) (flattened []map[string]interface{}) {
	for _, value := range values {
		flattened = append(flattened, map[string]interface{}{
			"synonyms": flattenStringList(value.Synonyms),
			"value":    aws.StringValue(value.Value),
		})
	}

	return
}

// Expects a slice of maps representing the Lex objects.
// The value passed into this function should have been run through the expandLexSet function.
// Example: []map[value: lilies synonyms:[]lirium]]
func expandLexEnumerationValues(rawValues []map[string]interface{}) (enums []*lexmodelbuildingservice.EnumerationValue) {
	for _, rawValue := range rawValues {
		enums = append(enums, &lexmodelbuildingservice.EnumerationValue{
			Synonyms: expandStringList(rawValue["synonyms"].([]interface{})),
			Value:    aws.String(rawValue["value"].(string)),
		})
	}

	return
}
