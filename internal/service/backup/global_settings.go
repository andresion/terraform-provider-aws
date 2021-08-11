package backup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/internal/client"
	"github.com/terraform-providers/terraform-provider-aws/internal/flex"
)

func ResourceGlobalSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsBackupGlobalSettingsUpdate,
		Update: resourceAwsBackupGlobalSettingsUpdate,
		Read:   resourceAwsBackupGlobalSettingsRead,
		Delete: schema.Noop,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"global_settings": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAwsBackupGlobalSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.AWSClient).BackupConn

	input := &backup.UpdateGlobalSettingsInput{
		GlobalSettings: flex.ExpandStringMap(d.Get("global_settings").(map[string]interface{})),
	}

	_, err := conn.UpdateGlobalSettings(input)
	if err != nil {
		return fmt.Errorf("error setting Backup Global Settings (%s): %w", meta.(*client.AWSClient).AccountID, err)
	}

	d.SetId(meta.(*client.AWSClient).AccountID)

	return resourceAwsBackupGlobalSettingsRead(d, meta)
}

func resourceAwsBackupGlobalSettingsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.AWSClient).BackupConn

	resp, err := conn.DescribeGlobalSettings(&backup.DescribeGlobalSettingsInput{})
	if err != nil {
		return fmt.Errorf("error reading Backup Global Settings (%s): %w", d.Id(), err)
	}

	if err := d.Set("global_settings", aws.StringValueMap(resp.GlobalSettings)); err != nil {
		return fmt.Errorf("error setting global_settings: %w", err)
	}

	return nil
}
