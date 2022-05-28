//go:build unit

package service_test

import (
	"github.com/emanuelefalzone/bitly/internal/application/service"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRandomKeyGenerator(t *testing.T) {
	// Build our needed test case data struct
	type testCase struct {
		key string
	}
	// Create new test cases
	testCases := []testCase{
		{
			key: "dd5w3b", // Depends on seed
		},
	}

	// Create new key generator with 0 seed
	generator := service.NewRandomKeyGenerator(0)
	for _, tc := range testCases {
		// Run Tests
		assert.Equal(t, tc.key, generator.NextKey())
	}
}
