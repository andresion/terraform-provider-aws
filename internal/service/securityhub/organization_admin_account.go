package securityhub

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/internal/client"
	"github.com/terraform-providers/terraform-provider-aws/internal/verify"
)

func ResourceOrganizationAdminAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceAwsSecurityHubOrganizationAdminAccountCreate,
		Read:   resourceAwsSecurityHubOrganizationAdminAccountRead,
		Delete: resourceAwsSecurityHubOrganizationAdminAccountDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"admin_account_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: verify.ValidAccountID,
			},
		},
	}
}

func resourceAwsSecurityHubOrganizationAdminAccountCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.AWSClient).SecurityHubConn

	adminAccountID := d.Get("admin_account_id").(string)

	input := &securityhub.EnableOrganizationAdminAccountInput{
		AdminAccountId: aws.String(adminAccountID),
	}

	_, err := conn.EnableOrganizationAdminAccount(input)

	if err != nil {
		return fmt.Errorf("error enabling Security Hub Organization Admin Account (%s): %w", adminAccountID, err)
	}

	d.SetId(adminAccountID)

	if _, err := waitAdminAccountEnabled(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for Security Hub Organization Admin Account (%s) to enable: %w", d.Id(), err)
	}

	return resourceAwsSecurityHubOrganizationAdminAccountRead(d, meta)
}

func resourceAwsSecurityHubOrganizationAdminAccountRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.AWSClient).SecurityHubConn

	adminAccount, err := findAdminAccount(conn, d.Id())

	if !d.IsNewResource() && tfawserr.ErrCodeEquals(err, securityhub.ErrCodeResourceNotFoundException) {
		log.Printf("[WARN] Security Hub Organization Admin Account (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading Security Hub Organization Admin Account (%s): %w", d.Id(), err)
	}

	if adminAccount == nil {
		if d.IsNewResource() {
			return fmt.Errorf("error reading Security Hub Organization Admin Account (%s): %w", d.Id(), err)
		}

		log.Printf("[WARN] Security Hub Organization Admin Account (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("admin_account_id", adminAccount.AccountId)

	return nil
}

func resourceAwsSecurityHubOrganizationAdminAccountDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*client.AWSClient).SecurityHubConn

	input := &securityhub.DisableOrganizationAdminAccountInput{
		AdminAccountId: aws.String(d.Id()),
	}

	_, err := conn.DisableOrganizationAdminAccount(input)

	if tfawserr.ErrCodeEquals(err, securityhub.ErrCodeResourceNotFoundException) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error disabling Security Hub Organization Admin Account (%s): %w", d.Id(), err)
	}

	if _, err := waitAdminAccountNotFound(conn, d.Id()); err != nil {
		return fmt.Errorf("error waiting for Security Hub Organization Admin Account (%s) to disable: %w", d.Id(), err)
	}

	return nil
}
