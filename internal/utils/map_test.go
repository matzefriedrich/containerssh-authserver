package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MergeMaps_add_missing(t *testing.T) {
	// Arrange
	m := map[string]string{"a": "b"}
	other := map[string]string{"c": "d"}

	// Act
	actual := MergeMaps(m, other, AddMissing)

	// Assert
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, "b", actual["a"])
	assert.Equal(t, "d", actual["c"])
}

func Test_MergeMaps_override(t *testing.T) {
	// Arrange
	m := map[string]string{"a": "b", "c": "d"}
	other := map[string]string{"a": "c"}

	// Act
	actual := MergeMaps(m, other, Override)

	// Assert
	assert.Equal(t, 2, len(actual))
	assert.Equal(t, "c", actual["a"])
	assert.Equal(t, "d", actual["c"])
}

func Test_MergeMaps_remove(t *testing.T) {
	// Arrange
	m := map[string]string{"a": "b"}
	other := map[string]string{"a": ""}

	// Act
	actual := MergeMaps(m, other, Remove)

	// Assert
	assert.Equal(t, 0, len(actual))
}
