package trie

import (
	"testing"
	"time"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestTrie(t *testing.T) {
	tests := map[string]struct {
		words      []string
		operations func(t *testing.T, trie *Trie)
	}{
		"empty should return nil for all": {
			words: []string{},
			operations: func(t *testing.T, trie *Trie) {
				t.Run("should be empty", func(t *testing.T) {
					assert.Empty(t, trie.GetCompletion(""))
					assert.False(t, trie.HasWord("some word"))
				})

				t.Run("should have value", func(t *testing.T) {
					trie.AddWords("test")
					expectedRes := []string{"test"}

					assert.Equal(t, expectedRes, trie.GetCompletion("t"))
					assert.Equal(t, expectedRes, trie.GetCompletion("te"))
					assert.Equal(t, expectedRes, trie.GetCompletion("tes"))
					assert.Equal(t, []string{}, trie.GetCompletion("a"))

					assert.True(t, trie.HasWord("test"))
					assert.Equal(t, expectedRes, trie.GetAllWords())
				})
			},
		},
		"has \"test\" value": {
			words: []string{"test"},
			operations: func(t *testing.T, trie *Trie) {
				t.Run("should have value", func(t *testing.T) {
					expectedRes := []string{"test"}

					assert.Equal(t, expectedRes, trie.GetCompletion("t"))
					assert.Equal(t, expectedRes, trie.GetCompletion("te"))
					assert.Equal(t, expectedRes, trie.GetCompletion("tes"))
					assert.Equal(t, []string{}, trie.GetCompletion("b"))

					assert.True(t, trie.HasWord("test"))
					assert.Equal(t, expectedRes, trie.GetAllWords())
				})
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			trie := NewTrie(tc.words...)
			tc.operations(t, trie)
		})
	}
}

const lenItems = 20

func BenchmarkTrie_Write(b *testing.B) {
	trie := NewTrie()
	for i := 0; i < b.N; i++ {
		ranStr := RandStringBytesMaskImprSrcUnsafe(lenItems)
		trie.AddWords(ranStr)
	}
}

const wordLens = 64

func BenchmarkTrie_Read(b *testing.B) {
	words := make([]string, 0, wordLens)
	for i := 0; i < wordLens; i++ {
		words = append(words, RandStringBytesMaskImprSrcUnsafe(i))
	}

	trie := NewTrie(words...)
	for j := 0; j < b.N; j++ {
		assert.True(b, trie.HasWord(words[j%wordLens]))
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(uint64(time.Now().UnixNano()))

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Uint64(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Uint64(), letterIdxMax
		}

		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
