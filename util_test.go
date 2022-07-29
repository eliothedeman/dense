package dense

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyEncoding(t *testing.T) {
	buff := make([]byte, 1024)
	assert.Equal(t, getKeyFunc[int]()(10, buff), []byte{0, 0, 0, 0, 0, 0, 0, 10})
	assert.Equal(t, getKeyFunc[string]()("a", buff), []byte{'a'})
	assert.Equal(t, getKeyFunc[string]()("abcd", buff), []byte{'a', 'b', 'c', 'd'})
	assert.Equal(t, getKeyFunc[[]byte]()([]byte{'b'}, buff), []byte{'b'})
	assert.Equal(t, getKeyFunc[uint]()(math.MaxUint, buff), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

}
