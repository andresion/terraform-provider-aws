package glue

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	mlTransformStatusUnknown   = "Unknown"
	registryStatusUnknown      = "Unknown"
	schemaStatusUnknown        = "Unknown"
	schemaVersionStatusUnknown = "Unknown"
	triggerStatusUnknown       = "Unknown"
)

// statusMLTransform fetches the MLTransform and its Status
func statusMLTransform(conn *glue.Glue, transformId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &glue.GetMLTransformInput{
			TransformId: aws.String(transformId),
		}

		output, err := conn.GetMLTransform(input)

		if err != nil {
			return nil, mlTransformStatusUnknown, err
		}

		if output == nil {
			return output, mlTransformStatusUnknown, nil
		}

		return output, aws.StringValue(output.Status), nil
	}
}

// statusRegistry fetches the Registry and its Status
func statusRegistry(conn *glue.Glue, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findRegistryByID(conn, id)
		if err != nil {
			return nil, registryStatusUnknown, err
		}

		if output == nil {
			return output, registryStatusUnknown, nil
		}

		return output, aws.StringValue(output.Status), nil
	}
}

// statusSchema fetches the Schema and its Status
func statusSchema(conn *glue.Glue, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findSchemaByID(conn, id)
		if err != nil {
			return nil, schemaStatusUnknown, err
		}

		if output == nil {
			return output, schemaStatusUnknown, nil
		}

		return output, aws.StringValue(output.SchemaStatus), nil
	}
}

// statusSchemaVersion fetches the Schema Version and its Status
func statusSchemaVersion(conn *glue.Glue, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		output, err := findSchemaVersionByID(conn, id)
		if err != nil {
			return nil, schemaVersionStatusUnknown, err
		}

		if output == nil {
			return output, schemaVersionStatusUnknown, nil
		}

		return output, aws.StringValue(output.Status), nil
	}
}

// statusTrigger fetches the Trigger and its Status
func statusTrigger(conn *glue.Glue, triggerName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		input := &glue.GetTriggerInput{
			Name: aws.String(triggerName),
		}

		output, err := conn.GetTrigger(input)

		if err != nil {
			return nil, triggerStatusUnknown, err
		}

		if output == nil {
			return output, triggerStatusUnknown, nil
		}

		return output, aws.StringValue(output.Trigger.State), nil
	}
}

func statusGlueDevEndpoint(conn *glue.Glue, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getDevEndpointInput := &glue.GetDevEndpointInput{
			EndpointName: aws.String(name),
		}
		endpoint, err := conn.GetDevEndpoint(getDevEndpointInput)
		if err != nil {
			if tfawserr.ErrCodeEquals(err, glue.ErrCodeEntityNotFoundException) {
				return nil, "", nil
			}

			return nil, "", err
		}

		if endpoint == nil || endpoint.DevEndpoint == nil {
			return nil, "", nil
		}

		if aws.StringValue(endpoint.DevEndpoint.Status) == "FAILED" && endpoint.DevEndpoint.FailureReason != nil {
			return endpoint, aws.StringValue(endpoint.DevEndpoint.Status), fmt.Errorf("%s", aws.StringValue(endpoint.DevEndpoint.FailureReason))
		}

		return endpoint, aws.StringValue(endpoint.DevEndpoint.Status), nil
	}
}
