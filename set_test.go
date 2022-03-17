package dense

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleInsert(t *testing.T) {

	s := NewSet[int](100)
	assert.False(t, s.Contains(100))
	s.Insert(100)
	assert.True(t, s.Contains(100))
}
