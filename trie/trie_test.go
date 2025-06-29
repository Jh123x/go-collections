package trie

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func BenchmarkTrie(b *testing.B) {

}
