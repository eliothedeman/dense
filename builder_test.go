package dense

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilderKeys(t *testing.T) {
	keys := []string{"hell", "hello"}
	b := builder[string]{}
	for _, k := range keys {
		b.Add([]byte(k), k)
	}
	m := b.Build()

	assert.Equal(t, m.MustGet([]byte("hell")), "hell")
	assert.Equal(t, m.MustGet([]byte("hello")), "hello")
}
