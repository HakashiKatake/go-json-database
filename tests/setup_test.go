package tests

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// TestGopterSetup verifies that gopter is properly configured
func TestGopterSetup(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("addition is commutative", prop.ForAll(
		func(a, b int) bool {
			return a+b == b+a
		},
		gen.Int(),
		gen.Int(),
	))

	properties.TestingRun(t)
}
