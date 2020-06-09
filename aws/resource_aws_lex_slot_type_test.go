package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAwsLexSlotType_basic(t *testing.T) {
	var v lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(rName, LexSlotTypeVersionLatest),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsLexSlotTypeBasicConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enumeration_value.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "value_selection_strategy", lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue),
					resource.TestCheckResourceAttrSet(resourceName, "checksum"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					testAccCheckResourceAttrRfc3339(resourceName, "created_date"),
					testAccCheckResourceAttrRfc3339(resourceName, "last_updated_date"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
		},
	})
}

func TestAccAwsLexSlotType_disappears(t *testing.T) {
	var v lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(rName, LexSlotTypeVersionLatest),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsLexSlotTypeBasicConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v),
					testAccCheckResourceDisappears(testAccProvider, resourceAwsLexSlotType(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccAwsLexSlotType_CreateVersion(t *testing.T) {
	var v1, v2, v3 lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(rName, LexSlotTypeVersionLatest),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsLexSlotTypeCreateVersionConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v1),
					testAccCheckAwsLexSlotTypeNotExists(rName, "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
			{
				Config: testAccAwsLexSlotTypeUpdateCreateVersionConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v2),
					testAccCheckAwsLexSlotTypeExistsWithVersion(rName, "1", &v3),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
		},
	})
}

func TestAccAwsLexSlotType_Description(t *testing.T) {
	var v1, v2 lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(rName, LexSlotTypeVersionLatest),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsLexSlotTypeBasicConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v1),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
			{
				Config: testAccAwsLexSlotTypeUpdateDescriptionConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v2),
					resource.TestCheckResourceAttr(resourceName, "description", "Types of flowers to pick up"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
		},
	})
}

func TestAccAwsLexSlotType_EnumerationValue(t *testing.T) {
	var v1, v2 lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(rName, LexSlotTypeVersionLatest),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsLexSlotTypeEnumerationValueConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v1),
					resource.TestCheckResourceAttr(resourceName, "enumeration_value.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
			{
				Config: testAccAwsLexSlotTypeUpdateEnumerationValueConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v2),
					resource.TestCheckResourceAttr(resourceName, "enumeration_value.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
		},
	})
}

func TestAccAwsLexSlotType_Name(t *testing.T) {
	var v1, v2 lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	testID1 := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testID2 := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	rName1 := "test_slot_type_" + testID1
	rName2 := "test_slot_type_" + testID2

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(rName1, LexSlotTypeVersionLatest),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsLexSlotTypeBasicConfig(rName1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName1, &v1),
					resource.TestCheckResourceAttr(resourceName, "name", rName1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
			{
				Config: testAccAwsLexSlotTypeBasicConfig(rName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName2, &v2),
					resource.TestCheckResourceAttr(resourceName, "name", rName2),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
		},
	})
}

func TestAccAwsLexSlotType_ValueSelectionStrategy(t *testing.T) {
	var v1, v2 lexmodelbuildingservice.GetSlotTypeOutput

	resourceName := "aws_lex_slot_type.test"
	rName := "test_slot_type_" + acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(rName, LexSlotTypeVersionLatest),
		Steps: []resource.TestStep{
			{
				Config: testAccAwsLexSlotTypeBasicConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v1),
					resource.TestCheckResourceAttr(resourceName, "value_selection_strategy", lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
			{
				Config: testAccAwsLexSlotTypeUpdateValueSelectionStrategyConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(rName, &v2),
					resource.TestCheckResourceAttr(resourceName, "value_selection_strategy", lexmodelbuildingservice.SlotValueSelectionStrategyTopResolution),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_version"},
			},
		},
	})
}

func testAccCheckAwsLexSlotTypeExists(slotTypeName string, res *lexmodelbuildingservice.GetSlotTypeOutput) resource.TestCheckFunc {
	return testAccCheckAwsLexSlotTypeExistsWithVersion(slotTypeName, LexSlotTypeVersionLatest, res)
}

func testAccCheckAwsLexSlotTypeExistsWithVersion(slotTypeName, slotTypeVersion string, res *lexmodelbuildingservice.GetSlotTypeOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).lexmodelconn

		out, err := conn.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
			Name:    aws.String(slotTypeName),
			Version: aws.String(slotTypeVersion),
		})
		if isAWSErr(err, lexmodelbuildingservice.ErrCodeNotFoundException, "") {
			return fmt.Errorf("error slot type %s not found, %s", slotTypeName, err)
		}
		if err != nil {
			return fmt.Errorf("error getting slot type %s: %s", slotTypeName, err)
		}

		*res = *out

		return nil
	}
}

func testAccCheckAwsLexSlotTypeNotExists(slotTypeName, slotTypeVersion string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).lexmodelconn

		_, err := conn.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
			Name:    aws.String(slotTypeName),
			Version: aws.String(slotTypeVersion),
		})
		if isAWSErr(err, lexmodelbuildingservice.ErrCodeNotFoundException, "") {
			return nil
		}
		if err != nil {
			return fmt.Errorf("error getting slot type %s: %s", slotTypeName, err)
		}

		return fmt.Errorf("error slot type exists %s", slotTypeName)
	}
}

func testAccAwsLexSlotTypeBasicConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_lex_slot_type" "test" {
  name = "%[1]s"
  enumeration_value {
    value = "roses"
  }
}
`, r)
}

func testAccAwsLexSlotTypeCreateVersionConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_lex_slot_type" "test" {
  name           = "%[1]s"
  create_version = false
  enumeration_value {
    value = "azaleas"
  }
}
`, r)
}

func testAccAwsLexSlotTypeUpdateCreateVersionConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_lex_slot_type" "test" {
  name           = "%[1]s"
  description    = "Types of flowers to pick up"
  create_version = true

  enumeration_value {
    value = "azaleas"
  }
}
`, r)
}

func testAccAwsLexSlotTypeUpdateDescriptionConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_lex_slot_type" "test" {
  name        = "%[1]s"
  description = "Types of flowers to pick up"
  enumeration_value {
    value = "chrysanthemums"
  }
}
`, r)
}

func testAccAwsLexSlotTypeEnumerationValueConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_lex_slot_type" "test" {
  name = "%[1]s"
  enumeration_value {
    value = "tulips"
    synonyms = [
      "Eduardoregelia",
      "Podonix",
    ]
  }
}
`, r)
}

func testAccAwsLexSlotTypeUpdateEnumerationValueConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_lex_slot_type" "test" {
  name = "%[1]s"

  enumeration_value {
    value = "lilies"
    synonyms = [
      "Lirium",
      "Martagon",
    ]
  }
  enumeration_value {
    value = "tulips"
    synonyms = [
      "Eduardoregelia",
      "Podonix",
    ]
  }
}
`, r)
}

func testAccAwsLexSlotTypeUpdateValueSelectionStrategyConfig(r string) string {
	return fmt.Sprintf(`
resource "aws_lex_slot_type" "test" {
  name                     = "%[1]s"
  value_selection_strategy = "TOP_RESOLUTION"
  enumeration_value {
    value = "dasies"
  }
}
`, r)
}
