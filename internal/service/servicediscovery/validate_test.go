package servicediscovery_test

import (
	"strings"
	"testing"
)

func TestValidNamespaceName(t *testing.T) {
	validNames := []string{
		"ValidName",
		"V_-.dN01e",
		"0",
		".",
		"-",
		"_",
		strings.Repeat("x", 1024),
	}
	for _, v := range validNames {
		_, errors := validateServiceDiscoveryNamespaceName(v, "name")
		if len(errors) != 0 {
			t.Fatalf("%q should be a valid namespace name: %q", v, errors)
		}
	}

	invalidNames := []string{
		"Inval:dName",
		"Invalid Name",
		"*",
		"",
		// length > 512
		strings.Repeat("x", 1025),
	}
	for _, v := range invalidNames {
		_, errors := validateServiceDiscoveryNamespaceName(v, "name")
		if len(errors) == 0 {
			t.Fatalf("%q should be an invalid namespace name", v)
		}
	}
}
