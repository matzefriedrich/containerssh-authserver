package utils

import (
	"gotest.tools/v3/assert"
	"testing"
)

func Test_IIf_returns_true_value_if_condition_is_true(t *testing.T) {
	// Arrange
	const expected = 1

	// Act
	actual := IIf(true, 1, 2)

	// Assert
	assert.Equal(t, expected, actual)
}

func Test_IIf_returns_false_value_if_condition_is_false(t *testing.T) {
	// Arrange
	const expected = 2

	// Act
	actual := IIf(false, 1, 2)

	// Assert
	assert.Equal(t, expected, actual)
}
