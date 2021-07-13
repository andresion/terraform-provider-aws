package acctest

//
// Was aws/provider_test.go.
// Removed many service-specific functions.
// Converted to public functions.
//

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	tfaws "github.com/terraform-providers/terraform-provider-aws/aws"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/envvar"
	tfnet "github.com/terraform-providers/terraform-provider-aws/aws/internal/net"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/provider/meta"
)

const (
	// Provider name for single configuration testing
	ProviderNameAws = "aws"

	// Provider name for alternate configuration testing
	ProviderNameAwsAlternate = "awsalternate"

	// Provider name for alternate account and alternate region configuration testing
	ProviderNameAwsAlternateAccountAlternateRegion = "awsalternateaccountalternateregion"

	// Provider name for alternate account and same region configuration testing
	ProviderNameAwsAlternateAccountSameRegion = "awsalternateaccountsameregion"

	// Provider name for same account and alternate region configuration testing
	ProviderNameAwsSameAccountAlternateRegion = "awssameaccountalternateregion"

	// Provider name for third configuration testing
	ProviderNameAwsThird = "awsthird"
)

const rfc3339RegexPattern = `^[0-9]{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9](\.[0-9]+)?([Zz]|([+-]([01][0-9]|2[0-3]):[0-5][0-9]))$`

// From validators.go.
const (
	awsAccountIDRegexpInternalPattern = `(aws|\d{12})`
	awsRegionRegexpInternalPattern    = `[a-z]{2}(-[a-z]+)+-\d`
)

// TestAccSkip implements a wrapper for (*testing.T).Skip() to prevent unused linting reports
//
// Reference: https://github.com/dominikh/go-tools/issues/633#issuecomment-606560616
var TestAccSkip = func(t *testing.T, message string) {
	t.Skip(message)
}

// testAccProviders is a static map containing only the main provider instance.
//
// Deprecated: Terraform Plugin SDK version 2 uses TestCase.ProviderFactories
// but supports this value in TestCase.Providers for backwards compatibility.
// In the future Providers: testAccProviders will be changed to
// ProviderFactories: testAccProviderFactories
var TestAccProviders map[string]*schema.Provider

// TestAccProviderFactories is a static map containing only the main provider instance
//
// Use other testAccProviderFactories functions, such as testAccProviderFactoriesAlternate,
// for tests requiring special provider configurations.
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// TestAccProvider is the "main" provider instance
//
// This Provider can be used in testing code for API calls without requiring
// the use of saving and referencing specific ProviderFactories instances.
//
// TestAccPreCheck(t) must be called before using this provider instance.
var TestAccProvider *schema.Provider

// testAccProviderConfigure ensures testAccProvider is only configured once
//
// The testAccPreCheck(t) function is invoked for every test and this prevents
// extraneous reconfiguration to the same values each time. However, this does
// not prevent reconfiguration that may happen should the address of
// testAccProvider be errantly reused in ProviderFactories.
var testAccProviderConfigure sync.Once

func init() {
	TestAccProvider = tfaws.Provider()

	TestAccProviders = map[string]*schema.Provider{
		ProviderNameAws: TestAccProvider,
	}

	// Always allocate a new provider instance each invocation, otherwise gRPC
	// ProviderConfigure() can overwrite configuration during concurrent testing.
	TestAccProviderFactories = map[string]func() (*schema.Provider, error){
		ProviderNameAws: func() (*schema.Provider, error) { return tfaws.Provider(), nil }, //nolint:unparam
	}
}

// TestAccProviderFactoriesInit creates ProviderFactories for the provider under testing.
func TestAccProviderFactoriesInit(providers *[]*schema.Provider, providerNames []string) map[string]func() (*schema.Provider, error) {
	var factories = make(map[string]func() (*schema.Provider, error), len(providerNames))

	for _, name := range providerNames {
		p := tfaws.Provider()

		factories[name] = func() (*schema.Provider, error) { //nolint:unparam
			return p, nil
		}

		if providers != nil {
			*providers = append(*providers, p)
		}
	}

	return factories
}

// TestAccProviderFactoriesInternal creates ProviderFactories for provider configuration testing
//
// This should only be used for TestAccAWSProvider_ tests which need to
// reference the provider instance itself. Other testing should use
// TestAccProviderFactories or other related functions.
func TestAccProviderFactoriesInternal(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error) {
	return TestAccProviderFactoriesInit(providers, []string{ProviderNameAws})
}

// TestAccProviderFactoriesAlternate creates ProviderFactories for cross-account and cross-region configurations
//
// For cross-region testing: Typically paired with TestAccMultipleRegionPreCheck and TestAccAlternateRegionProviderConfig.
//
// For cross-account testing: Typically paired with TestAccAlternateAccountPreCheck and TestAccAlternateAccountProviderConfig.
func TestAccProviderFactoriesAlternate(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error) {
	return TestAccProviderFactoriesInit(providers, []string{
		ProviderNameAws,
		ProviderNameAwsAlternate,
	})
}

// TestAccProviderFactoriesAlternateAccountAndAlternateRegion creates ProviderFactories for cross-account and cross-region configurations
//
// Usage typically paired with TestAccMultipleRegionPreCheck, TestAccAlternateAccountPreCheck,
// and TestAccAlternateAccountAndAlternateRegionProviderConfig.
func TestAccProviderFactoriesAlternateAccountAndAlternateRegion(providers *[]*schema.Provider) map[string]func() (*schema.Provider, error) {
	return TestAccProviderFactoriesInit(providers, []string{
		ProviderNameAws,
		ProviderNameAwsAlternateAccountAlternateRegion,
		ProviderNameAwsAlternateAccountSameRegion,
		ProviderNameAwsSameAccountAlternateRegion,
	})
}

// TestAccProviderFactoriesMultipleRegion creates ProviderFactories for the number of region configurations
//
// Usage typically paired with TestAccMultipleRegionPreCheck and TestAccMultipleRegionProviderConfig.
func TestAccProviderFactoriesMultipleRegion(providers *[]*schema.Provider, regions int) map[string]func() (*schema.Provider, error) {
	providerNames := []string{
		ProviderNameAws,
		ProviderNameAwsAlternate,
	}

	if regions >= 3 {
		providerNames = append(providerNames, ProviderNameAwsThird)
	}

	return TestAccProviderFactoriesInit(providers, providerNames)
}

// TestAccPreCheck verifies and sets required provider testing configuration
//
// This PreCheck function should be present in every acceptance test. It allows
// test configurations to omit a provider configuration with region and ensures
// testing functions that attempt to call AWS APIs are previously configured.
//
// These verifications and configuration are preferred at this level to prevent
// provider developers from experiencing less clear errors for every test.
func TestAccPreCheck(t *testing.T) {
	// Since we are outside the scope of the Terraform configuration we must
	// call Configure() to properly initialize the provider configuration.
	testAccProviderConfigure.Do(func() {
		envvar.TestFailIfAllEmpty(t, []string{envvar.AwsProfile, envvar.AwsAccessKeyId, envvar.AwsContainerCredentialsFullUri}, "credentials for running acceptance testing")

		if os.Getenv(envvar.AwsAccessKeyId) != "" {
			envvar.TestFailIfEmpty(t, envvar.AwsSecretAccessKey, "static credentials value when using "+envvar.AwsAccessKeyId)
		}

		// Setting the AWS_DEFAULT_REGION environment variable here allows all tests to omit
		// a provider configuration with a region. This defaults to us-west-2 for provider
		// developer simplicity and has been in the codebase for a very long time.
		//
		// This handling must be preserved until either:
		//   * AWS_DEFAULT_REGION is required and checked above (should mention us-west-2 default)
		//   * Region is automatically handled via shared AWS configuration file and still verified
		region := TestAccGetRegion()
		os.Setenv(envvar.AwsDefaultRegion, region)

		err := TestAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
		if err != nil {
			t.Fatal(err)
		}
	})
}

// TestAccAwsProviderAccountID returns the account ID of an AWS provider
func TestAccAwsProviderAccountID(provider *schema.Provider) string {
	if provider == nil {
		log.Print("[DEBUG] Unable to read account ID from test provider: empty provider")
		return ""
	}
	if provider.Meta() == nil {
		log.Print("[DEBUG] Unable to read account ID from test provider: unconfigured provider")
		return ""
	}
	client, ok := provider.Meta().(meta.Meta)
	if !ok {
		log.Print("[DEBUG] Unable to read account ID from test provider: non-AWS or unconfigured AWS provider")
		return ""
	}
	return client.GetAccountID()
}

// TestAccCheckResourceAttrAccountID ensures the Terraform state exactly matches the account ID
func TestAccCheckResourceAttrAccountID(resourceName, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(resourceName, attributeName, TestAccGetAccountID())(s)
	}
}

// TestAccCheckResourceAttrRegionalARN ensures the Terraform state exactly matches a formatted ARN with region
func TestAccCheckResourceAttrRegionalARN(resourceName, attributeName, arnService, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributeValue := arn.ARN{
			AccountID: TestAccGetAccountID(),
			Partition: TestAccGetPartition(),
			Region:    TestAccGetRegion(),
			Resource:  arnResource,
			Service:   arnService,
		}.String()
		return resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue)(s)
	}
}

// TestAccCheckResourceAttrRegionalARNNoAccount ensures the Terraform state exactly matches a formatted ARN with region but without account ID
func TestAccCheckResourceAttrRegionalARNNoAccount(resourceName, attributeName, arnService, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributeValue := arn.ARN{
			Partition: TestAccGetPartition(),
			Region:    TestAccGetRegion(),
			Resource:  arnResource,
			Service:   arnService,
		}.String()
		return resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue)(s)
	}
}

// TestAccCheckResourceAttrRegionalARNAccountID ensures the Terraform state exactly matches a formatted ARN with region and specific account ID
func TestAccCheckResourceAttrRegionalARNAccountID(resourceName, attributeName, arnService, accountID, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributeValue := arn.ARN{
			AccountID: accountID,
			Partition: TestAccGetPartition(),
			Region:    TestAccGetRegion(),
			Resource:  arnResource,
			Service:   arnService,
		}.String()
		return resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue)(s)
	}
}

// TestAccMatchResourceAttrAccountID ensures the Terraform state regexp matches an account ID
func TestAccMatchResourceAttrAccountID(resourceName, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestMatchResourceAttr(resourceName, attributeName, regexp.MustCompile(`^\d{12}$`))(s)
	}
}

// TestAccMatchResourceAttrRegionalARN ensures the Terraform state regexp matches a formatted ARN with region
func TestAccMatchResourceAttrRegionalARN(resourceName, attributeName, arnService string, arnResourceRegexp *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		arnRegexp := arn.ARN{
			AccountID: TestAccGetAccountID(),
			Partition: TestAccGetPartition(),
			Region:    TestAccGetRegion(),
			Resource:  arnResourceRegexp.String(),
			Service:   arnService,
		}.String()

		attributeMatch, err := regexp.Compile(arnRegexp)

		if err != nil {
			return fmt.Errorf("Unable to compile ARN regexp (%s): %w", arnRegexp, err)
		}

		return resource.TestMatchResourceAttr(resourceName, attributeName, attributeMatch)(s)
	}
}

// TestAccMatchResourceAttrRegionalARNNoAccount ensures the Terraform state regexp matches a formatted ARN with region but without account ID
func TestAccMatchResourceAttrRegionalARNNoAccount(resourceName, attributeName, arnService string, arnResourceRegexp *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		arnRegexp := arn.ARN{
			Partition: TestAccGetPartition(),
			Region:    TestAccGetRegion(),
			Resource:  arnResourceRegexp.String(),
			Service:   arnService,
		}.String()

		attributeMatch, err := regexp.Compile(arnRegexp)

		if err != nil {
			return fmt.Errorf("Unable to compile ARN regexp (%s): %s", arnRegexp, err)
		}

		return resource.TestMatchResourceAttr(resourceName, attributeName, attributeMatch)(s)
	}
}

// TestAccMatchResourceAttrRegionalARNAccountID ensures the Terraform state regexp matches a formatted ARN with region and specific account ID
func TestAccMatchResourceAttrRegionalARNAccountID(resourceName, attributeName, arnService, accountID string, arnResourceRegexp *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		arnRegexp := arn.ARN{
			AccountID: accountID,
			Partition: TestAccGetPartition(),
			Region:    TestAccGetRegion(),
			Resource:  arnResourceRegexp.String(),
			Service:   arnService,
		}.String()

		attributeMatch, err := regexp.Compile(arnRegexp)

		if err != nil {
			return fmt.Errorf("Unable to compile ARN regexp (%s): %w", arnRegexp, err)
		}

		return resource.TestMatchResourceAttr(resourceName, attributeName, attributeMatch)(s)
	}
}

// TestAccMatchResourceAttrRegionalHostname ensures the Terraform state regexp matches a formatted DNS hostname with region and partition DNS suffix
func TestAccMatchResourceAttrRegionalHostname(resourceName, attributeName, serviceName string, hostnamePrefixRegexp *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		hostnameRegexpPattern := fmt.Sprintf("%s\\.%s\\.%s\\.%s$", hostnamePrefixRegexp.String(), serviceName, TestAccGetRegion(), TestAccGetPartitionDNSSuffix())

		hostnameRegexp, err := regexp.Compile(hostnameRegexpPattern)

		if err != nil {
			return fmt.Errorf("Unable to compile hostname regexp (%s): %w", hostnameRegexp, err)
		}

		return resource.TestMatchResourceAttr(resourceName, attributeName, hostnameRegexp)(s)
	}
}

// TestAccCheckResourceAttrGlobalARN ensures the Terraform state exactly matches a formatted ARN without region
func TestAccCheckResourceAttrGlobalARN(resourceName, attributeName, arnService, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributeValue := arn.ARN{
			AccountID: TestAccGetAccountID(),
			Partition: TestAccGetPartition(),
			Resource:  arnResource,
			Service:   arnService,
		}.String()
		return resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue)(s)
	}
}

// TestAccCheckResourceAttrGlobalARNNoAccount ensures the Terraform state exactly matches a formatted ARN without region or account ID
func TestAccCheckResourceAttrGlobalARNNoAccount(resourceName, attributeName, arnService, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributeValue := arn.ARN{
			Partition: TestAccGetPartition(),
			Resource:  arnResource,
			Service:   arnService,
		}.String()
		return resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue)(s)
	}
}

// TestAccCheckResourceAttrGlobalARNAccountID ensures the Terraform state exactly matches a formatted ARN without region and with specific account ID
func TestAccCheckResourceAttrGlobalARNAccountID(resourceName, attributeName, accountID, arnService, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributeValue := arn.ARN{
			AccountID: accountID,
			Partition: TestAccGetPartition(),
			Resource:  arnResource,
			Service:   arnService,
		}.String()
		return resource.TestCheckResourceAttr(resourceName, attributeName, attributeValue)(s)
	}
}

// TestAccMatchResourceAttrGlobalARN ensures the Terraform state regexp matches a formatted ARN without region
func TestAccMatchResourceAttrGlobalARN(resourceName, attributeName, arnService string, arnResourceRegexp *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		arnRegexp := arn.ARN{
			AccountID: TestAccGetAccountID(),
			Partition: TestAccGetPartition(),
			Resource:  arnResourceRegexp.String(),
			Service:   arnService,
		}.String()

		attributeMatch, err := regexp.Compile(arnRegexp)

		if err != nil {
			return fmt.Errorf("Unable to compile ARN regexp (%s): %w", arnRegexp, err)
		}

		return resource.TestMatchResourceAttr(resourceName, attributeName, attributeMatch)(s)
	}
}

// TestAccCheckResourceAttrRegionalARNIgnoreRegionAndAccount ensures the Terraform state exactly matches a formatted ARN with region without specifying the region or account
func TestAccCheckResourceAttrRegionalARNIgnoreRegionAndAccount(resourceName, attributeName, arnService, arnResource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		arnRegexp := arn.ARN{
			AccountID: awsAccountIDRegexpInternalPattern,
			Partition: TestAccGetPartition(),
			Region:    awsRegionRegexpInternalPattern,
			Resource:  arnResource,
			Service:   arnService,
		}.String()

		attributeMatch, err := regexp.Compile(arnRegexp)

		if err != nil {
			return fmt.Errorf("Unable to compile ARN regexp (%s): %w", arnRegexp, err)
		}

		return resource.TestMatchResourceAttr(resourceName, attributeName, attributeMatch)(s)
	}
}

// TestAccMatchResourceAttrGlobalARNNoAccount ensures the Terraform state regexp matches a formatted ARN without region or account ID
func TestAccMatchResourceAttrGlobalARNNoAccount(resourceName, attributeName, arnService string, arnResourceRegexp *regexp.Regexp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		arnRegexp := arn.ARN{
			Partition: TestAccGetPartition(),
			Resource:  arnResourceRegexp.String(),
			Service:   arnService,
		}.String()

		attributeMatch, err := regexp.Compile(arnRegexp)

		if err != nil {
			return fmt.Errorf("Unable to compile ARN regexp (%s): %s", arnRegexp, err)
		}

		return resource.TestMatchResourceAttr(resourceName, attributeName, attributeMatch)(s)
	}
}

// TestAccCheckResourceAttrRfc3339 ensures the Terraform state matches a RFC3339 value
// This TestCheckFunc will likely be moved to the Terraform Plugin SDK in the future.
func TestAccCheckResourceAttrRfc3339(resourceName, attributeName string) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(resourceName, attributeName, regexp.MustCompile(rfc3339RegexPattern))
}

// TestAccCheckListHasSomeElementAttrPair is a TestCheckFunc which validates that the collection on the left has an element with an attribute value
// matching the value on the left
// Based on TestCheckResourceAttrPair from the Terraform SDK testing framework
func TestAccCheckListHasSomeElementAttrPair(nameFirst string, resourceAttr string, elementAttr string, nameSecond string, keySecond string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		isFirst, err := primaryInstanceState(s, nameFirst)
		if err != nil {
			return err
		}

		isSecond, err := primaryInstanceState(s, nameSecond)
		if err != nil {
			return err
		}

		vSecond, ok := isSecond.Attributes[keySecond]
		if !ok {
			return fmt.Errorf("%s: No attribute %q found", nameSecond, keySecond)
		} else if vSecond == "" {
			return fmt.Errorf("%s: No value was set on attribute %q", nameSecond, keySecond)
		}

		attrsFirst := make([]string, 0, len(isFirst.Attributes))
		attrMatcher := regexp.MustCompile(fmt.Sprintf("%s\\.\\d+\\.%s", resourceAttr, elementAttr))
		for k := range isFirst.Attributes {
			if attrMatcher.MatchString(k) {
				attrsFirst = append(attrsFirst, k)
			}
		}

		found := false
		for _, attrName := range attrsFirst {
			vFirst := isFirst.Attributes[attrName]
			if vFirst == vSecond {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("%s: No element of %q found with attribute %q matching value %q set on %q of %s", nameFirst, resourceAttr, elementAttr, vSecond, keySecond, nameSecond)
		}

		return nil
	}
}

// TestAccCheckResourceAttrEquivalentJSON is a TestCheckFunc that compares a JSON value with an expected value. Both JSON
// values are normalized before being compared.
func TestAccCheckResourceAttrEquivalentJSON(resourceName, attributeName, expectedJSON string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		is, err := primaryInstanceState(s, resourceName)
		if err != nil {
			return err
		}

		v, ok := is.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("%s: No attribute %q found", resourceName, attributeName)
		}

		vNormal, err := structure.NormalizeJsonString(v)
		if err != nil {
			return fmt.Errorf("%s: Error normalizing JSON in %q: %w", resourceName, attributeName, err)
		}

		expectedNormal, err := structure.NormalizeJsonString(expectedJSON)
		if err != nil {
			return fmt.Errorf("Error normalizing expected JSON: %w", err)
		}

		if vNormal != expectedNormal {
			return fmt.Errorf("%s: Attribute %q expected\n%s\ngot\n%s", resourceName, attributeName, expectedJSON, v)
		}
		return nil
	}
}

// Copied and inlined from the SDK testing code
func primaryInstanceState(s *terraform.State, name string) (*terraform.InstanceState, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("Not found: %s", name)
	}

	is := rs.Primary
	if is == nil {
		return nil, fmt.Errorf("No primary instance: %s", name)
	}

	return is, nil
}

// TestAccGetAccountID returns the account ID of testAccProvider
// Must be used within a resource.TestCheckFunc
func TestAccGetAccountID() string {
	return TestAccAwsProviderAccountID(TestAccProvider)
}

func TestAccGetRegion() string {
	return envvar.GetWithDefault(envvar.AwsDefaultRegion, endpoints.UsWest2RegionID)
}

func TestAccGetAlternateRegion() string {
	return envvar.GetWithDefault(envvar.AwsAlternateRegion, endpoints.UsEast1RegionID)
}

func TestAccGetThirdRegion() string {
	return envvar.GetWithDefault(envvar.AwsThirdRegion, endpoints.UsEast2RegionID)
}

func TestAccGetPartition() string {
	if partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), TestAccGetRegion()); ok {
		return partition.ID()
	}
	return "aws"
}

func TestAccGetPartitionDNSSuffix() string {
	if partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), TestAccGetRegion()); ok {
		return partition.DNSSuffix()
	}
	return "amazonaws.com"
}

func TestAccGetPartitionReverseDNSPrefix() string {
	if partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), TestAccGetRegion()); ok {
		return tfnet.ReverseDns(partition.DNSSuffix())
	}

	return "com.amazonaws"
}

func TestAccGetAlternateRegionPartition() string {
	if partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), TestAccGetAlternateRegion()); ok {
		return partition.ID()
	}
	return "aws"
}

func TestAccGetThirdRegionPartition() string {
	if partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), TestAccGetThirdRegion()); ok {
		return partition.ID()
	}
	return "aws"
}

func TestAccAlternateAccountPreCheck(t *testing.T) {
	envvar.TestSkipIfAllEmpty(t, []string{envvar.AwsAlternateProfile, envvar.AwsAlternateAccessKeyId}, "credentials for running acceptance testing in alternate AWS account")

	if os.Getenv(envvar.AwsAlternateAccessKeyId) != "" {
		envvar.TestSkipIfEmpty(t, envvar.AwsAlternateSecretAccessKey, "static credentials value when using "+envvar.AwsAlternateAccessKeyId)
	}
}

func TestAccPartitionHasServicePreCheck(serviceId string, t *testing.T) {
	if partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), TestAccGetRegion()); ok {
		if _, ok := partition.Services()[serviceId]; !ok {
			t.Skip(fmt.Sprintf("skipping tests; partition %s does not support %s service", partition.ID(), serviceId))
		}
	}
}

func TestAccMultipleRegionPreCheck(t *testing.T, regions int) {
	if TestAccGetRegion() == TestAccGetAlternateRegion() {
		t.Fatalf("%s and %s must be set to different values for acceptance tests", envvar.AwsDefaultRegion, envvar.AwsAlternateRegion)
	}

	if TestAccGetPartition() != TestAccGetAlternateRegionPartition() {
		t.Fatalf("%s partition (%s) does not match %s partition (%s)", envvar.AwsAlternateRegion, TestAccGetAlternateRegionPartition(), envvar.AwsDefaultRegion, TestAccGetPartition())
	}

	if regions >= 3 {
		if TestAccGetRegion() == TestAccGetThirdRegion() {
			t.Fatalf("%s and %s must be set to different values for acceptance tests", envvar.AwsDefaultRegion, envvar.AwsThirdRegion)
		}

		if TestAccGetAlternateRegion() == TestAccGetThirdRegion() {
			t.Fatalf("%s and %s must be set to different values for acceptance tests", envvar.AwsAlternateRegion, envvar.AwsThirdRegion)
		}

		if TestAccGetPartition() != TestAccGetThirdRegionPartition() {
			t.Fatalf("%s partition (%s) does not match %s partition (%s)", envvar.AwsThirdRegion, TestAccGetThirdRegionPartition(), envvar.AwsDefaultRegion, TestAccGetPartition())
		}
	}

	if partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), TestAccGetRegion()); ok {
		if len(partition.Regions()) < regions {
			t.Skipf("skipping tests; partition includes %d regions, %d expected", len(partition.Regions()), regions)
		}
	}
}

// TestAccRegionPreCheck checks that the test region is the specified region.
func TestAccRegionPreCheck(t *testing.T, region string) {
	if TestAccGetRegion() != region {
		t.Skipf("skipping tests; %s (%s) does not equal %s", envvar.AwsDefaultRegion, TestAccGetRegion(), region)
	}
}

// TestAccPartitionPreCheck checks that the test partition is the specified partition.
func TestAccPartitionPreCheck(partition string, t *testing.T) {
	if TestAccGetPartition() != partition {
		t.Skipf("skipping tests; current partition (%s) does not equal %s", TestAccGetPartition(), partition)
	}
}

func TestAccAlternateAccountProviderConfig() string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "awsalternate" {
  access_key = %[1]q
  profile    = %[2]q
  secret_key = %[3]q
}
`, os.Getenv(envvar.AwsAlternateAccessKeyId), os.Getenv(envvar.AwsAlternateProfile), os.Getenv(envvar.AwsAlternateSecretAccessKey))
}

func TestAccAlternateAccountAlternateRegionProviderConfig() string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "awsalternate" {
  access_key = %[1]q
  profile    = %[2]q
  region     = %[3]q
  secret_key = %[4]q
}
`, os.Getenv(envvar.AwsAlternateAccessKeyId), os.Getenv(envvar.AwsAlternateProfile), TestAccGetAlternateRegion(), os.Getenv(envvar.AwsAlternateSecretAccessKey))
}

// When testing needs to distinguish a second region and second account in the same region
// e.g. cross-region functionality with RAM shared subnets
func TestAccAlternateAccountAndAlternateRegionProviderConfig() string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "awsalternateaccountalternateregion" {
  access_key = %[1]q
  profile    = %[2]q
  region     = %[3]q
  secret_key = %[4]q
}

provider "awsalternateaccountsameregion" {
  access_key = %[1]q
  profile    = %[2]q
  secret_key = %[4]q
}

provider "awssameaccountalternateregion" {
  region = %[3]q
}
`, os.Getenv(envvar.AwsAlternateAccessKeyId), os.Getenv(envvar.AwsAlternateProfile), TestAccGetAlternateRegion(), os.Getenv(envvar.AwsAlternateSecretAccessKey))
}

// Deprecated: Use TestAccMultipleRegionProviderConfig instead
func TestAccAlternateRegionProviderConfig() string {
	return TestAccNamedRegionalProviderConfig(ProviderNameAwsAlternate, TestAccGetAlternateRegion())
}

func TestAccMultipleRegionProviderConfig(regions int) string {
	var config strings.Builder

	config.WriteString(TestAccNamedRegionalProviderConfig(ProviderNameAwsAlternate, TestAccGetAlternateRegion()))

	if regions >= 3 {
		config.WriteString(TestAccNamedRegionalProviderConfig(ProviderNameAwsThird, TestAccGetThirdRegion()))
	}

	return config.String()
}

func TestAccProviderConfigDefaultAndIgnoreTagsKeyPrefixes1(key1, value1, keyPrefix1 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  default_tags {
    tags = {
      %q = %q
    }
  }
  ignore_tags {
    key_prefixes = [%q]
  }
}
`, key1, value1, keyPrefix1)
}

func TestAccProviderConfigDefaultAndIgnoreTagsKeys1(key1, value1 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  default_tags {
    tags = {
      %[1]q = %q
    }
  }
  ignore_tags {
    keys = [%[1]q]
  }
}
`, key1, value1)
}

func TestAccProviderConfigIgnoreTagsKeyPrefixes1(keyPrefix1 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  ignore_tags {
    key_prefixes = [%[1]q]
  }
}
`, keyPrefix1)
}

func TestAccProviderConfigIgnoreTagsKeys1(key1 string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  ignore_tags {
    keys = [%[1]q]
  }
}
`, key1)
}

// TestAccNamedRegionalProviderConfig creates a new provider named configuration with a region.
//
// This can be used to build multiple provider configuration testing.
func TestAccNamedRegionalProviderConfig(providerName string, region string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider %[1]q {
  region = %[2]q
}
`, providerName, region)
}

// TestAccRegionalProviderConfig creates a new provider configuration with a region.
//
// This can only be used for single provider configuration testing as it
// overwrites the "aws" provider configuration.
func TestAccRegionalProviderConfig(region string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  region = %[1]q
}
`, region)
}

func TestAccAwsRegionProviderFunc(region string, providers *[]*schema.Provider) func() *schema.Provider {
	return func() *schema.Provider {
		if region == "" {
			log.Println("[DEBUG] No region given")
			return nil
		}
		if providers == nil {
			log.Println("[DEBUG] No providers given")
			return nil
		}

		log.Printf("[DEBUG] Checking providers for AWS region: %s", region)
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider == nil || provider.Meta() == nil {
				log.Printf("[DEBUG] Skipping empty provider")
				continue
			}

			// Ignore if Meta is not AWSClient, this will happen for other providers
			client, ok := provider.Meta().(meta.Meta)
			if !ok {
				log.Printf("[DEBUG] Skipping non-AWS provider")
				continue
			}

			clientRegion := client.GetRegion()
			log.Printf("[DEBUG] Checking AWS provider region %q against %q", clientRegion, region)
			if clientRegion == region {
				log.Printf("[DEBUG] Found AWS provider with region: %s", region)
				return provider
			}
		}

		log.Printf("[DEBUG] No suitable provider found for %q in %d providers", region, len(*providers))
		return nil
	}
}

func TestAccDeleteResource(resource *schema.Resource, d *schema.ResourceData, meta interface{}) error {
	if resource.DeleteContext != nil || resource.DeleteWithoutTimeout != nil {
		var diags diag.Diagnostics

		if resource.DeleteContext != nil {
			diags = resource.DeleteContext(context.Background(), d, meta)
		} else {
			diags = resource.DeleteWithoutTimeout(context.Background(), d, meta)
		}

		for i := range diags {
			if diags[i].Severity == diag.Error {
				return fmt.Errorf("error deleting resource: %s", diags[i].Summary)
			}
		}

		return nil
	}

	return resource.Delete(d, meta)
}

func TestAccCheckResourceDisappears(provider *schema.Provider, resourceFunc func(interface{}) *schema.Resource, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if resourceState.Primary.ID == "" {
			return fmt.Errorf("resource ID missing: %s", resourceName)
		}

		meta := provider.Meta()
		resource := resourceFunc(meta)

		return TestAccDeleteResource(resource, resource.Data(resourceState.Primary), meta)
	}
}

func TestAccCheckWithProviders(f func(*terraform.State, *schema.Provider) error, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		numberOfProviders := len(*providers)
		for i, provider := range *providers {
			if provider.Meta() == nil {
				log.Printf("[DEBUG] Skipping empty provider %d (total: %d)", i, numberOfProviders)
				continue
			}
			log.Printf("[DEBUG] Calling check with provider %d (total: %d)", i, numberOfProviders)
			if err := f(s, provider); err != nil {
				return err
			}
		}
		return nil
	}
}

// TestAccErrorCheckSkipMessagesContaining skips tests based on error messages that indicate unsupported features
func TestAccErrorCheckSkipMessagesContaining(t *testing.T, messages ...string) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}

		for _, message := range messages {
			errorMessage := err.Error()
			if strings.Contains(errorMessage, message) {
				t.Skipf("skipping test for %s/%s: %s", TestAccGetPartition(), TestAccGetRegion(), errorMessage)
			}
		}

		return err
	}
}

type ServiceErrorCheckFunc func(*testing.T) resource.ErrorCheckFunc

var serviceErrorCheckFuncs map[string]ServiceErrorCheckFunc

func RegisterServiceErrorCheckFunc(endpointID string, f ServiceErrorCheckFunc) {
	if serviceErrorCheckFuncs == nil {
		serviceErrorCheckFuncs = make(map[string]ServiceErrorCheckFunc)
	}

	if _, ok := serviceErrorCheckFuncs[endpointID]; ok {
		// already registered
		panic(fmt.Sprintf("Cannot re-register a service! ServiceErrorCheckFunc exists for %s", endpointID)) //lintignore:R009
	}

	serviceErrorCheckFuncs[endpointID] = f
}

func TestAccErrorCheck(t *testing.T, endpointIDs ...string) resource.ErrorCheckFunc {
	return func(err error) error {
		if err == nil {
			return nil
		}

		for _, endpointID := range endpointIDs {
			if f, ok := serviceErrorCheckFuncs[endpointID]; ok {
				ef := f(t)
				err = ef(err)
			}

			if err == nil {
				break
			}
		}

		if TestAccErrorCheckCommon(err) {
			t.Skipf("skipping test for %s/%s: %s", TestAccGetPartition(), TestAccGetRegion(), err.Error())
		}

		return err
	}
}

// NOTE: This function cannot use the standard tfawserr helpers
// as it is receiving error strings from the SDK testing framework,
// not actual error types from the resource logic.
func TestAccErrorCheckCommon(err error) bool {
	if strings.Contains(err.Error(), "is not supported in this") {
		return true
	}

	if strings.Contains(err.Error(), "is currently not supported") {
		return true
	}

	if strings.Contains(err.Error(), "InvalidAction") {
		return true
	}

	if strings.Contains(err.Error(), "Unknown operation") {
		return true
	}

	if strings.Contains(err.Error(), "UnknownOperationException") {
		return true
	}

	if strings.Contains(err.Error(), "UnsupportedOperation") {
		return true
	}

	return false
}

// Check service API call error for reasons to skip acceptance testing
// These include missing API endpoints and unsupported API calls
func TestAccPreCheckSkipError(err error) bool {
	// GovCloud has endpoints that respond with (no message provided after the error code):
	// AccessDeniedException:
	// Ignore these API endpoints that exist but are not officially enabled
	if tfawserr.ErrCodeEquals(err, "AccessDeniedException") {
		return true
	}
	// Ignore missing API endpoints
	if tfawserr.ErrMessageContains(err, "RequestError", "send request failed") {
		return true
	}
	// Ignore unsupported API calls
	if tfawserr.ErrCodeEquals(err, "UnknownOperationException") {
		return true
	}
	if tfawserr.ErrCodeEquals(err, "UnsupportedOperation") {
		return true
	}
	if tfawserr.ErrMessageContains(err, "InvalidInputException", "Unknown operation") {
		return true
	}
	if tfawserr.ErrMessageContains(err, "InvalidAction", "is not valid") {
		return true
	}
	if tfawserr.ErrMessageContains(err, "InvalidAction", "Unavailable Operation") {
		return true
	}
	return false
}

func TestAccAWSProviderConfigDefaultTags_Tags0() string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		`
provider "aws" {
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`)
}

func TestAccAWSProviderConfigDefaultTags_Tags1(tag1, value1 string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  default_tags {
    tags = {
      %q = %q
    }
  }

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`, tag1, value1))
}

func TestAccAWSProviderConfigDefaultTags_Tags2(tag1, value1, tag2, value2 string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  default_tags {
    tags = {
      %q = %q
      %q = %q
    }
  }

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`, tag1, value1, tag2, value2))
}

func TestAccAWSProviderConfigDefaultTagsEmptyConfigurationBlock() string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		`
provider "aws" {
  default_tags {}

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`)
}

func TestAccAWSProviderConfigDefaultAndIgnoreTagsEmptyConfigurationBlock() string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		`
provider "aws" {
  default_tags {}
  ignore_tags {}

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`)
}

func TestAccAWSProviderConfigEndpoints(endpoints string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    %[1]s
  }
}
`, endpoints))
}

func TestAccAWSProviderConfigIgnoreTagsEmptyConfigurationBlock() string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		`
provider "aws" {
  ignore_tags {}

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`)
}

func TestAccAWSProviderConfigIgnoreTagsKeyPrefixes0() string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		`
provider "aws" {
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`)
}

func TestAccAWSProviderConfigIgnoreTagsKeyPrefixes1(tagPrefix1 string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  ignore_tags {
    key_prefixes = [%[1]q]
  }

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`, tagPrefix1))
}

func TestAccAWSProviderConfigIgnoreTagsKeyPrefixes2(tagPrefix1, tagPrefix2 string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  ignore_tags {
    key_prefixes = [%[1]q, %[2]q]
  }

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`, tagPrefix1, tagPrefix2))
}

func TestAccAWSProviderConfigIgnoreTagsKeys0() string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		`
provider "aws" {
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`)
}

func TestAccAWSProviderConfigIgnoreTagsKeys1(tag1 string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  ignore_tags {
    keys = [%[1]q]
  }

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`, tag1))
}

func TestAccAWSProviderConfigIgnoreTagsKeys2(tag1, tag2 string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  ignore_tags {
    keys = [%[1]q, %[2]q]
  }

  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`, tag1, tag2))
}

func TestAccAWSProviderConfigRegion(region string) string {
	//lintignore:AT004
	return ComposeConfig(
		TestAccProviderConfigBase,
		fmt.Sprintf(`
provider "aws" {
  region                      = %[1]q
  skip_credentials_validation = true
  skip_get_ec2_platforms      = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true
}
`, region))
}

func TestAccAssumeRoleARNPreCheck(t *testing.T) {
	envvar.TestSkipIfEmpty(t, envvar.TfAccAssumeRoleArn, "Amazon Resource Name (ARN) of existing IAM Role to assume for testing restricted permissions")
}

func TestAccProviderConfigAssumeRolePolicy(policy string) string {
	//lintignore:AT004
	return fmt.Sprintf(`
provider "aws" {
  assume_role {
    role_arn = %q
    policy   = %q
  }
}
`, os.Getenv(envvar.TfAccAssumeRoleArn), policy)
}

const TestAccCheckAWSProviderConfigAssumeRoleEmpty = `
provider "aws" {
  assume_role {
  }
}

data "aws_caller_identity" "current" {}
` //lintignore:AT004

const TestAccProviderConfigBase = `
data "aws_partition" "provider_test" {}

# Required to initialize the provider
data "aws_arn" "test" {
  arn = "arn:${data.aws_partition.provider_test.partition}:s3:::test"
}
`

func TestCheckResourceAttrIsSortedCsv(resourceName, attributeName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		is, err := primaryInstanceState(s, resourceName)
		if err != nil {
			return err
		}

		v, ok := is.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("%s: No attribute %q found", resourceName, attributeName)
		}

		splitV := strings.Split(v, ",")
		if !sort.StringsAreSorted(splitV) {
			return fmt.Errorf("%s: Expected attribute %q to be sorted, got %q", resourceName, attributeName, v)
		}

		return nil
	}
}

// composeConfig can be called to concatenate multiple strings to build test configurations
func ComposeConfig(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}
