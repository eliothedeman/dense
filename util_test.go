package dense

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyEncoding(t *testing.T) {
	buff := make([]byte, 1024)
	assert.Equal(t, getKeyFunc[int]()(10, buff), bitvec{0, 0, 0, 0, 0, 0, 0, 10})
	assert.Equal(t, getKeyFunc[string]()("a", buff), bitvec{'a'})
	assert.Equal(t, getKeyFunc[string]()("abcd", buff), bitvec{'a', 'b', 'c', 'd'})
	assert.Equal(t, getKeyFunc[[]byte]()([]byte{'b'}, buff), bitvec{'b'})
	assert.Equal(t, getKeyFunc[uint]()(math.MaxUint, buff), bitvec{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
}
