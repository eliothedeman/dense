package dense

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBitvecGetSet(t *testing.T) {
	b := make(bitvec, 2)
	b.set(2, 2)
	assert.Equal(t, uint8(2), b.at(2))
	b.set(1, 3)
	assert.Equal(t, uint8(3), b.at(1))
}
