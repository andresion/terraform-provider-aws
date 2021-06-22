package dataconv

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ExpandStringMap expands a map of string to interface to a map of string to *string.
func ExpandStringMap(m map[string]interface{}) map[string]*string {
	stringMap := make(map[string]*string, len(m))
	for k, v := range m {
		stringMap[k] = aws.String(v.(string))
	}
	return stringMap
}

// ExpandStringSet takes the result of schema.Set of strings and returns a []*string.
func ExpandStringSet(configured *schema.Set) []*string {
	return ExpandStringList(configured.List()) // nosemgrep: helper-schema-Set-extraneous-expandStringList-with-List
}

// ExpandStringList takes the result of flatmap.Expand for an array of strings and returns a []*string.
func ExpandStringList(configured []interface{}) []*string {
	vs := make([]*string, 0, len(configured))
	for _, v := range configured {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, aws.String(v.(string)))
		}
	}
	return vs
}
