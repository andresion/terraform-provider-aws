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

func TestAccLexSlotType_basic(t *testing.T) {
	resourceName := "aws_lex_slot_type.test"
	testID := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testSlotTypeID := "test_slot_type_" + testID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(testSlotTypeID, "$LATEST"),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeBasicConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enumeration_value.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", testSlotTypeID),
					resource.TestCheckResourceAttr(resourceName, "value_selection_strategy", lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue),
					resource.TestCheckResourceAttrSet(resourceName, "checksum"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsLexSlotType_CreateVersion(t *testing.T) {
	resourceName := "aws_lex_slot_type.test"
	testID := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testSlotTypeID := "test_slot_type_" + testID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(testSlotTypeID, "$LATEST"),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeCreateVersionConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					testAccCheckAwsLexSlotTypeNotExists(testSlotTypeID, "1"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
				// Cannot verify import for create_version other than true because we don't have
				// that info at import, it is not returned from the AWS API.
				// ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeUpdateCreateVersionConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsLexSlotType_Description(t *testing.T) {
	resourceName := "aws_lex_slot_type.test"
	testID := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testSlotTypeID := "test_slot_type_" + testID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(testSlotTypeID, "$LATEST"),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeBasicConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeUpdateDescriptionConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "description", "Types of flowers to pick up"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsLexSlotType_EnumerationValue(t *testing.T) {
	resourceName := "aws_lex_slot_type.test"
	testID := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testSlotTypeID := "test_slot_type_" + testID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(testSlotTypeID, "$LATEST"),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeEnumerationValueConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "enumeration_value.#", "1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeUpdateEnumerationValueConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "enumeration_value.#", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsLexSlotType_Name(t *testing.T) {
	resourceName := "aws_lex_slot_type.test"
	testID1 := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testID2 := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testSlotTypeID1 := "test_slot_type_" + testID1
	testSlotTypeID2 := "test_slot_type_" + testID2

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(testSlotTypeID1, "$LATEST"),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeBasicConfig, testSlotTypeID1),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID1, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "name", testSlotTypeID1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeBasicConfig, testSlotTypeID2),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID2, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "name", testSlotTypeID2),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsLexSlotType_ValueSelectionStrategy(t *testing.T) {
	resourceName := "aws_lex_slot_type.test"
	testID := acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
	testSlotTypeID := "test_slot_type_" + testID

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAwsLexSlotTypeNotExists(testSlotTypeID, "$LATEST"),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeBasicConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "value_selection_strategy", lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccAwsLexSlotTypeUpdateValueSelectionStrategyConfig, testSlotTypeID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAwsLexSlotTypeExists(testSlotTypeID, "$LATEST"),
					resource.TestCheckResourceAttr(resourceName, "value_selection_strategy", lexmodelbuildingservice.SlotValueSelectionStrategyTopResolution),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAwsLexSlotTypeExists(slotTypeName, slotTypeVersion string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*AWSClient).lexmodelconn

		_, err := conn.GetSlotType(&lexmodelbuildingservice.GetSlotTypeInput{
			Name:    aws.String(slotTypeName),
			Version: aws.String(slotTypeVersion),
		})
		if isAWSErr(err, lexmodelbuildingservice.ErrCodeNotFoundException, "") {
			return fmt.Errorf("error slot type %s not found, %s", slotTypeName, err)
		}
		if err != nil {
			return fmt.Errorf("error getting slot type %s: %s", slotTypeName, err)
		}

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

const testAccAwsLexSlotTypeBasicConfig = `
resource "aws_lex_slot_type" "test" {
  name = "%s"
}
`

const testAccAwsLexSlotTypeCreateVersionConfig = `
resource "aws_lex_slot_type" "test" {
  name           = "%s"
  create_version = false
}
`
const testAccAwsLexSlotTypeUpdateCreateVersionConfig = `
resource "aws_lex_slot_type" "test" {
  name           = "%s"
  description    = "Types of flowers to pick up"
  create_version = true
}
`

const testAccAwsLexSlotTypeUpdateDescriptionConfig = `
resource "aws_lex_slot_type" "test" {
  description = "Types of flowers to pick up"
  name        = "%s"
}
`

const testAccAwsLexSlotTypeEnumerationValueConfig = `
resource "aws_lex_slot_type" "test" {
  enumeration_value {
    synonyms = [
      "Eduardoregelia",
      "Podonix",
    ]

    value = "tulips"
  }

  name = "%s"
}
`

const testAccAwsLexSlotTypeUpdateEnumerationValueConfig = `
resource "aws_lex_slot_type" "test" {
  enumeration_value {
    synonyms = [
      "Lirium",
      "Martagon",
    ]

    value = "lilies"
  }

  enumeration_value {
    synonyms = [
      "Eduardoregelia",
      "Podonix",
    ]

    value = "tulips"
  }

  name = "%s"
}
`

const testAccAwsLexSlotTypeUpdateValueSelectionStrategyConfig = `
resource "aws_lex_slot_type" "test" {
  name                     = "%s"
  value_selection_strategy = "TOP_RESOLUTION"
}
`
