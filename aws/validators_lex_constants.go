package aws

import (
	"github.com/aws/aws-sdk-go/service/lexmodelbuildingservice"
)

// Amazon Lex Resource Constants. Data models are documented here
// https://docs.aws.amazon.com/lex/latest/dg/API_Types_Amazon_Lex_Model_Building_Service.html

const (

	// General

	lexNameMinLength = 1
	lexNameMaxLength = 100
	lexNameRegex     = "^([A-Za-z]_?)+$"

	lexVersionMinLength = 1
	lexVersionMaxLength = 64
	lexVersionRegex     = "\\$LATEST|[0-9]+"
	lexVersionLatest    = "$LATEST"
	lexVersionDefault   = "$LATEST"

	lexDescriptionMinLength = 0
	lexDescriptionMaxLength = 200

	// Slot Type

	lexSlotTypeMinLength                     = 1
	lexSlotTypeMaxLength                     = 100
	lexSlotTypeRegex                         = "^((AMAZON\\.)_?|[A-Za-z]_?)+"
	lexSlotTypeValueSelectionStrategyDefault = lexmodelbuildingservice.SlotValueSelectionStrategyOriginalValue

	// Enumeration Value

	lexEnumerationValuesMin             = 1
	lexEnumerationValuesMax             = 10000
	lexEnumerationValueSynonymMinLength = 1
	lexEnumerationValueSynonymMaxLength = 140
	lexEnumerationValueMinLength        = 1
	lexEnumerationValueMaxLength        = 140
)
