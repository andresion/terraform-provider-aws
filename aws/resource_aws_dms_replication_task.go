package aws

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	dms "github.com/aws/aws-sdk-go/service/databasemigrationservice"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/databasemigrationservice/finder"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/service/databasemigrationservice/waiter"
)

func resourceAwsDmsReplicationTask() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsDmsReplicationTaskCreate,
		Read:   resourceAwsDmsReplicationTaskRead,
		Update: resourceAwsDmsReplicationTaskUpdate,
		Delete: resourceAwsDmsReplicationTaskDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"cdc_start_time": {
				Type:     schema.TypeString,
				Optional: true,
				// Requires a Unix timestamp in seconds. Example 1484346880
			},
			"migration_type": {
				Type:     schema.TypeString,
				Required: true,
				// Per User Guide: "You can't change the migration type of a task."
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(dms.MigrationTypeValue_Values(), false),
			},
			"replication_instance_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArn,
			},
			"replication_task_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replication_task_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDmsReplicationTaskId,
			},
			"replication_task_settings": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
			},
			"source_endpoint_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArn,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_mappings": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: suppressEquivalentJsonDiffs,
			},
			"tags":     tagsSchema(),
			"tags_all": tagsSchemaComputed(),
			"target_endpoint_arn": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateArn,
			},
		},

		CustomizeDiff: SetTagsDiff,
	}
}

func resourceAwsDmsReplicationTaskCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dmsconn
	defaultTagsConfig := meta.(*AWSClient).DefaultTagsConfig
	tags := defaultTagsConfig.MergeTags(keyvaluetags.New(d.Get("tags").(map[string]interface{})))

	taskId := d.Get("replication_task_id").(string)

	request := &dms.CreateReplicationTaskInput{
		MigrationType:             aws.String(d.Get("migration_type").(string)),
		ReplicationInstanceArn:    aws.String(d.Get("replication_instance_arn").(string)),
		ReplicationTaskIdentifier: aws.String(taskId),
		SourceEndpointArn:         aws.String(d.Get("source_endpoint_arn").(string)),
		TableMappings:             aws.String(d.Get("table_mappings").(string)),
		Tags:                      tags.IgnoreAws().DatabasemigrationserviceTags(),
		TargetEndpointArn:         aws.String(d.Get("target_endpoint_arn").(string)),
	}

	if v, ok := d.GetOk("cdc_start_time"); ok {
		seconds, err := strconv.ParseInt(v.(string), 10, 64)
		if err != nil {
			return fmt.Errorf("DMS create replication task. Invalid CDC Unix timestamp: %s", err)
		}
		request.CdcStartTime = aws.Time(time.Unix(seconds, 0))
	}

	if v, ok := d.GetOk("replication_task_settings"); ok {
		request.ReplicationTaskSettings = aws.String(v.(string))
	}

	_, err := conn.CreateReplicationTask(request)

	if err != nil {
		return fmt.Errorf("error creating DMS replication task (%s): %w", taskId, err)
	}

	d.SetId(taskId)

	if err := waiter.ReplicationTaskReady(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for DMS replication task (%s) creation: %w", d.Id(), err)
	}

	return resourceAwsDmsReplicationTaskRead(d, meta)
}

func resourceAwsDmsReplicationTaskRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dmsconn
	defaultTagsConfig := meta.(*AWSClient).DefaultTagsConfig
	ignoreTagsConfig := meta.(*AWSClient).IgnoreTagsConfig

	task, err := finder.ReplicationTask(conn, d.Id())

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, dms.ErrCodeResourceNotFoundFault) {
		log.Printf("[WARN] DMS Replication Task (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading DMS replication task (%s): %w", d.Id(), err)
	}

	if task == nil {
		if d.IsNewResource() {
			return fmt.Errorf("error reading DMS replication task (%s): empty output after creation", d.Id())
		}
		log.Printf("[WARN] DMS Replication Task (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("migration_type", task.MigrationType)
	d.Set("replication_instance_arn", task.ReplicationInstanceArn)
	d.Set("replication_task_arn", task.ReplicationTaskArn)
	d.Set("replication_task_id", task.ReplicationTaskIdentifier)
	d.Set("source_endpoint_arn", task.SourceEndpointArn)
	d.Set("status", task.Status)
	d.Set("table_mappings", task.TableMappings)
	d.Set("target_endpoint_arn", task.TargetEndpointArn)

	settings, err := dmsReplicationTaskRemoveReadOnlySettings(aws.StringValue(task.ReplicationTaskSettings))
	if err != nil {
		return fmt.Errorf("error setting replication_task_settings: %w", err)
	}

	d.Set("replication_task_settings", settings)

	arn := aws.StringValue(task.ReplicationTaskArn)

	tags, err := keyvaluetags.DatabasemigrationserviceListTags(conn, arn)

	if err != nil {
		return fmt.Errorf("error listing tags for DMS Replication Task (%s): %s", arn, err)
	}

	tags = tags.IgnoreAws().IgnoreConfig(ignoreTagsConfig)

	//lintignore:AWSR002
	if err := d.Set("tags", tags.RemoveDefaultConfig(defaultTagsConfig).Map()); err != nil {
		return fmt.Errorf("error setting tags: %w", err)
	}

	if err := d.Set("tags_all", tags.Map()); err != nil {
		return fmt.Errorf("error setting tags_all: %w", err)
	}

	return nil
}

func resourceAwsDmsReplicationTaskUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dmsconn

	if d.HasChangesExcept("tags", "tags_all") {
		request := &dms.ModifyReplicationTaskInput{
			ReplicationTaskArn: aws.String(d.Get("replication_task_arn").(string)),
		}

		if d.HasChange("cdc_start_time") {
			seconds, err := strconv.ParseInt(d.Get("cdc_start_time").(string), 10, 64)
			if err != nil {
				return fmt.Errorf("DMS update replication task. Invalid CRC Unix timestamp: %s", err)
			}
			request.CdcStartTime = aws.Time(time.Unix(seconds, 0))
		}

		if d.HasChange("replication_task_settings") {
			request.ReplicationTaskSettings = aws.String(d.Get("replication_task_settings").(string))
		}

		if d.HasChange("table_mappings") {
			request.TableMappings = aws.String(d.Get("table_mappings").(string))
		}

		if d.Get("status").(string) == waiter.ReplicationTaskStatusRunning {
			input := &dms.StopReplicationTaskInput{
				ReplicationTaskArn: aws.String(d.Get("replication_task_arn").(string)),
			}

			_, err := conn.StopReplicationTask(input)

			if err != nil {
				return fmt.Errorf("error stopping DMS replication task (%s) before modifying: %w", d.Id(), err)
			}

			if err := waiter.ReplicationTaskStopped(conn, d.Id()); err != nil {
				return fmt.Errorf("error waiting for DMS replication task (%s) to be stopped: %w", d.Id(), err)
			}
		}

		_, err := conn.ModifyReplicationTask(request)

		if err != nil {
			return fmt.Errorf("error modifying DMS replication task (%s): %w", d.Id(), err)
		}

		if err := waiter.ReplicationTaskReady(conn, d.Id()); err != nil {
			return fmt.Errorf("error waiting for DMS replication task (%s) to be modified: %w", d.Id(), err)
		}
	}

	if d.HasChange("tags_all") {
		arn := d.Get("replication_task_arn").(string)
		o, n := d.GetChange("tags_all")

		if err := keyvaluetags.DatabasemigrationserviceUpdateTags(conn, arn, o, n); err != nil {
			return fmt.Errorf("error updating DMS Replication Task (%s) tags: %s", arn, err)
		}
	}

	return resourceAwsDmsReplicationTaskRead(d, meta)
}

func resourceAwsDmsReplicationTaskDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).dmsconn

	request := &dms.DeleteReplicationTaskInput{
		ReplicationTaskArn: aws.String(d.Get("replication_task_arn").(string)),
	}

	_, err := conn.DeleteReplicationTask(request)

	if tfawserr.ErrCodeEquals(err, dms.ErrCodeResourceNotFoundFault) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting DMS replication task (%s): %w", d.Id(), err)
	}

	if err := waiter.ReplicationTaskDeleted(conn, d.Id()); err != nil {
		if tfawserr.ErrCodeEquals(err, dms.ErrCodeResourceNotFoundFault) {
			return nil
		}

		return fmt.Errorf("error waiting for DMS replication task (%s) deletion: %w", d.Id(), err)
	}

	return nil
}

func dmsReplicationTaskRemoveReadOnlySettings(settings string) (*string, error) {
	var settingsData map[string]interface{}
	if err := json.Unmarshal([]byte(settings), &settingsData); err != nil {
		return nil, err
	}

	controlTablesSettings, ok := settingsData["ControlTablesSettings"].(map[string]interface{})
	if ok {
		delete(controlTablesSettings, "historyTimeslotInMinutes")
	}

	logging, ok := settingsData["Logging"].(map[string]interface{})
	if ok {
		delete(logging, "CloudWatchLogGroup")
		delete(logging, "CloudWatchLogStream")
	}

	cleanedSettings, err := json.Marshal(settingsData)
	if err != nil {
		return nil, err
	}

	cleanedSettingsString := string(cleanedSettings)
	return &cleanedSettingsString, nil
}
