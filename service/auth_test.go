package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

	// assert equality
	assert.Equal(t, []byte("STeZg1g5IEwyGlD/5fiBjrJ+WtXDlU2SxKMWlJuwAAM="), TokenKey, "they should be equal")

	// assert inequality
	assert.NotEqual(t, 123, 456, "they should not be equal")

}
