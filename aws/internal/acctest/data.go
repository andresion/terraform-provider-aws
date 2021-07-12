package acctest

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

type TestData struct{}

// RandomName returns a new random name with the standard prefix `tf-acc-test`.
func (td *TestData) RandomName() string {
	return acctest.RandomWithPrefix("tf-acc-test")
}

// NewTestData returns a new TestData structure.
func NewTestData(t *testing.T) TestData {
	data := TestData{}

	return data
}
