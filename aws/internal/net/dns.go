package net

import (
	"strings"
)

// ReverseDns switches a DNS hostname to reverse DNS and vice-versa.
func ReverseDns(hostname string) string {
	parts := strings.Split(hostname, ".")

	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, ".")
}
