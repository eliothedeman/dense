package dense

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkInsert(b *testing.B) {
	var buff bytes.Buffer
	b.Run("trie", func(b *testing.B) {
		x := NewTrie[int]()
		for i := 0; i < b.N; i++ {
			fmt.Fprint(&buff, i)
			x.Insert(buff.Bytes(), i)
			buff.Reset()
		}
	})
	b.Run("hash", func(b *testing.B) {
		x := make(map[string]int)
		for i := 0; i < b.N; i++ {
			fmt.Fprint(&buff, i)
			x[buff.String()] = i
			buff.Reset()
		}
	})
}

func TestInsertGet(t *testing.T) {
	x := NewTrie[int]()
	x.Insert([]byte("hello"), 10)
	assert.Equal(t, x.MustGet([]byte("hello")), 10)
	assert.Panics(t, func() {
		x.MustGet([]byte("nope"))
	})
}

func Test_bitsAtDepth(t *testing.T) {
	type args struct {
		data  []byte
		depth id
	}
	tests := []struct {
		name string
		args args
		want byte
	}{
		{
			"empty",
			args{
				[]byte{0, 0},
				0,
			},
			0,
		},
		{
			"at byte one",
			args{
				[]byte{0b01},
				0b01,
			},
			0,
		},
		{
			"at byte one offset 1",
			args{
				[]byte{0b0100},
				0b01,
			},
			0b01,
		},
		{
			"16",
			args{
				[]byte{0b10000},
				2,
			},
			0b01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bitsAtDepth(tt.args.data, tt.args.depth); got != tt.want {
				t.Errorf("bitsAtDepth() = %b, want %b", got, tt.want)
			}
		})
	}
}
