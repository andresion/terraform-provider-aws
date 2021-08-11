package ec2_test

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aws/internal/acctest"
	"github.com/terraform-providers/terraform-provider-aws/internal/client"
	"github.com/terraform-providers/terraform-provider-aws/internal/provider"
)

const (
	// EC2-Classic region testing environment variable name
	Ec2ClassicRegionEnvVar = "AWS_EC2_CLASSIC_REGION"
)

// testAccProviderEc2Classic is the EC2-Classic provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// testAccEC2ClassicPreCheck(t) must be called before using this provider instance.
var testAccProviderEc2Classic *schema.Provider

// testAccProviderEc2ClassicConfigure ensures the provider is only configured once
var testAccProviderEc2ClassicConfigure sync.Once

// testAccEC2ClassicPreCheck verifies AWS credentials and that EC2-Classic is supported
func testAccEC2ClassicPreCheck(t *testing.T) {
	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderEc2ClassicConfigure.Do(func() {
		testAccProviderEc2Classic = provider.Provider()

		config := map[string]interface{}{
			"region": testAccGetEc2ClassicRegion(),
		}

		err := testAccProviderEc2Classic.Configure(context.Background(), terraform.NewResourceConfigRaw(config))

		if err != nil {
			t.Fatal(err)
		}
	})

	client := testAccProviderEc2Classic.Meta().(*client.AWSClient)
	platforms := client.SupportedPlatforms
	region := client.Region
	if !client.HasEC2Classic(platforms) {
		t.Skipf("this test can only run in EC2-Classic, platforms available in %s: %q", region, platforms)
	}
}

// testAccEc2ClassicRegionProviderConfig is the Terraform provider configuration for EC2-Classic region testing
//
// Testing EC2-Classic assumes no other provider configurations are necessary
// and overwrites the "aws" provider configuration.
func testAccEc2ClassicRegionProviderConfig() string {
	return acctest.ConfigRegionalProvider(testAccGetEc2ClassicRegion())
}

// testAccGetEc2ClassicRegion returns the EC2-Classic region for testing
func testAccGetEc2ClassicRegion() string {
	v := os.Getenv(Ec2ClassicRegionEnvVar)

	if v != "" {
		return v
	}

	if acctest.Partition() == endpoints.AwsPartitionID {
		return endpoints.UsEast1RegionID
	}

	return acctest.Region()
}

// testAccCheckResourceAttrRegionalARNEc2Classic ensures the Terraform state exactly matches a formatted ARN with EC2-Classic region
func testAccCheckResourceAttrRegionalARNEc2Classic(resourceName, attributeName, arnService, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributeValue := arn.ARN{
			AccountID: acctest.AccountID(),
			Partition: acctest.Partition(),
			Region:    testAccGetEc2ClassicRegion(),
			Resource:  arnResource,
			Service:   arnService,
		}.String()
		return resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue)(s)
	}
}
