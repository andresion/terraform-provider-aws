package aws

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// Many of the Lex resources require complex nested objects. Terraform maps only support simple key
// value pairs and not complex or mixed types. That is why these resources are defined using the
// schema.TypeList and a max of 1 item instead of the schema.TypeMap.

func flattenLexEnumerationValues(values []*lexmodelbuildingservice.EnumerationValue) (flattened []map[string]interface{}) {
	for _, value := range values {
		flattened = append(flattened, map[string]interface{}{
			"synonyms": flattenStringList(value.Synonyms),
			"value":    aws.StringValue(value.Value),
		})
	}

	return
}

func expandLexEnumerationValues(rawValues []interface{}) []*lexmodelbuildingservice.EnumerationValue {
	enums := make([]*lexmodelbuildingservice.EnumerationValue, 0, len(rawValues))
	for _, rawValue := range rawValues {
		value, ok := rawValue.(map[string]interface{})
		if !ok {
			continue
		}

		enums = append(enums, &lexmodelbuildingservice.EnumerationValue{
			Synonyms: expandStringList(value["synonyms"].([]interface{})),
			Value:    aws.String(value["value"].(string)),
		})
	}
	return enums
}

func resourceAwsLexSlotType() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsLexSlotTypeCreate,
		Read:   resourceAwsLexSlotTypeRead,
		Update: resourceAwsLexSlotTypeUpdate,
		Delete: resourceAwsLexSlotTypeDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				// The version is not required for import but it is required for the get request.
				d.Set("version", "$LATEST")
				return []*schema.ResourceData{d}, nil
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute),
			Update: schema.DefaultTimeout(time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"checksum": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 200),
			},
			"enumeration_value": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 10000,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"synonyms": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringLenBetween(1, 140),
							},
						},
						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 140),
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(regexp.MustCompile(`^((AMAZON\.)_?|[A-Za-z]_?)+`), ""),
				),
			},
			"value_selection_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue,
				ValidateFunc: validation.StringInSlice([]string{
					lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue,
					lexmodelbuildingservice.SlotValueSelectionStrategyTopResolution,
				}, false),
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "$LATEST",
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile(`\$LATEST|[0-9]+`), ""),
				),
			},
		},
	}
}

func resourceAwsLexSlotTypeCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelconn
	name := d.Get("name").(string)

	input := &lexmodelbuildingservice.PutSlotTypeInput{
		Name:                   aws.String(name),
		ValueSelectionStrategy: aws.String(d.Get("value_selection_strategy").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("enumeration_value"); ok {
		input.EnumerationValues = expandLexEnumerationValues(v.(*schema.Set).List())
	}

	err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		output, err := conn.PutSlotType(input)

		if isAWSErr(err, lexmodelbuildingservice.ErrCodeConflictException, "") {
			input.Checksum = output.Checksum
			return resource.RetryableError(fmt.Errorf("%q: slot type still creating", d.Id()))
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error creating slot type %s: %s", name, err)
	}

	d.SetId(name)

	return resourceAwsLexSlotTypeRead(d, meta)
}

func resourceAwsLexSlotTypeRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelconn

	resp, err := conn.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
		Name:    aws.String(d.Id()),
		Version: aws.String(d.Get("version").(string)),
	})
	if isAWSErr(err, lexmodelbuildingservice.ErrCodeNotFoundException, "") {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error getting slot type %s: %s", d.Id(), err)
	}

	d.Set("checksum", resp.Checksum)
	d.Set("description", resp.Description)
	d.Set("name", resp.Name)
	d.Set("value_selection_strategy", resp.ValueSelectionStrategy)
	d.Set("version", resp.Version)

	if resp.EnumerationValues != nil {
		d.Set("enumeration_value", flattenLexEnumerationValues(resp.EnumerationValues))
	}

	return nil
}

func resourceAwsLexSlotTypeUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelconn

	input := &lexmodelbuildingservice.PutSlotTypeInput{
		Name:                   aws.String(d.Id()),
		Checksum:               aws.String(d.Get("checksum").(string)),
		CreateVersion:          aws.Bool(true),
		ValueSelectionStrategy: aws.String(d.Get("value_selection_strategy").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = aws.String(v.(string))
	}

	if v, ok := d.GetOk("enumeration_value"); ok {
		input.EnumerationValues = expandLexEnumerationValues(v.(*schema.Set).List())
	}

	err := resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err := conn.PutSlotType(input)

		if isAWSErr(err, lexmodelbuildingservice.ErrCodeConflictException, "") {
			return resource.RetryableError(fmt.Errorf("%q: slot type still updating", d.Id()))
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error updating slot type %s: %s", d.Id(), err)
	}

	return resourceAwsLexSlotTypeRead(d, meta)
}

func resourceAwsLexSlotTypeDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).lexmodelconn

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.DeleteSlotType(&lexmodelbuildingservice.DeleteSlotTypeInput{
			Name: aws.String(d.Id()),
		})

		if isAWSErr(err, lexmodelbuildingservice.ErrCodeConflictException, "") {
			return resource.RetryableError(fmt.Errorf("%q: slot type still deleting", d.Id()))
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("error deleting slot type %s: %s", d.Id(), err)
	}

	return nil
}
