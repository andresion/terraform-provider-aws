package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAwsLexSlotType_basic(t *testing.T) {
	resourceName := "aws_lex_slot_type.test"
	dataSourceName := "data." + resourceName
	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAwsLexSlotTypeConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "checksum", resourceName, "checksum"),
					resource.TestCheckResourceAttrPair(dataSourceName, "description", resourceName, "description"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enumeration_value.#", resourceName, "enumeration_value.#"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "value_selection_strategy", resourceName, "value_selection_strategy"),
					resource.TestCheckResourceAttrPair(dataSourceName, "version", resourceName, "version"),
					resource.TestCheckResourceAttrPair(dataSourceName, "created_date", resourceName, "created_date"),
					resource.TestCheckResourceAttrPair(dataSourceName, "last_updated_date", resourceName, "last_updated_date"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsLexSlotType_Version(t *testing.T) {
	var v1, v2, v3 lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	initialDataSourceName := "data.aws_lex_slot_type.initial"
	updatedDataSourceName := "data.aws_lex_slot_type.updated"

	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	initialVersion := "1"
	updatedVersion := "2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAwsLexSlotTypeConfig_VersionSetup(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExistsWithVersion(rName, initialVersion, &v1),
				),
			},
			{
				Config: testDataSourceAwsLexSlotTypeConfig_Version(rName, initialVersion),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExistsWithVersion(rName, initialVersion, &v2),
					resource.TestCheckResourceAttr(initialDataSourceName, "version", initialVersion),
					resource.TestCheckResourceAttrPair(initialDataSourceName, "created_date", resourceName, "created_date"),
					// The following value is for the $LATEST version, regardless of the verion retrieved
					resource.TestCheckResourceAttrPair(initialDataSourceName, "last_updated_date", resourceName, "last_updated_date"),

					testAccCheckAwsLexSlotTypeExistsWithVersion(rName, updatedVersion, &v3),
					resource.TestCheckResourceAttr(updatedDataSourceName, "version", updatedVersion),
					resource.TestCheckResourceAttrPair(updatedDataSourceName, "created_date", resourceName, "created_date"),
					resource.TestCheckResourceAttrPair(updatedDataSourceName, "last_updated_date", resourceName, "last_updated_date"),
				),
			},
		},
	})
}

func testDataSourceAwsLexSlotTypeConfig(r string) string {
	return composeConfig(
		testAccAwsLexSlotTypeBasicConfig(r), `
data "aws_lex_slot_type" "test" {
  name = aws_lex_slot_type.test.name
}
`)
}

func testDataSourceAwsLexSlotTypeConfig_VersionSetup(r string) string {
	return testAccAwsLexSlotTypeEnumerationValueConfig(r)
}

func testDataSourceAwsLexSlotTypeConfig_Version(r, v string) string {
	return composeConfig(
		testAccAwsLexSlotTypeUpdateEnumerationValueConfig(r),
		fmt.Sprintf(`
data "aws_lex_slot_type" "initial" {
  name    = aws_lex_slot_type.test.name
  version = %[1]q
}

data "aws_lex_slot_type" "updated" {
  name    = aws_lex_slot_type.test.name
}
`, v))
}
