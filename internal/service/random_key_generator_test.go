//go:build unit

package service_test

import (
	"github.com/emanuelefalzone/bitly/internal/service"
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestRandomKeyGenerator_NextKey(t *testing.T) {
	// Build our needed testcase data struct
	type testCase struct {
		location string
		key      string
	}
	// Create new test cases
	testCases := []testCase{
		{
			location: "http://www.google.com",
			key:      "dd5w3b",
		},
	}

	generator := service.NewRandomKeyGenerator(0)
	for _, tc := range testCases {
		// Run Tests
		assert.Equal(t, tc.key, generator.NextKey(tc.location))
	}
}
