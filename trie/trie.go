package trie

import "fmt"

type Trie struct {
	head *Node
}

func NewTrie(data ...string) *Trie {
	trie := &Trie{head: NewNode()}
	trie.AddWords(data...)

	return trie
}

func (t *Trie) AddWords(words ...string) {
	for _, word := range words {
		t.head.AddWord(word)
	}
}

func (t *Trie) HasWord(word string) bool {
	return t.head.HasWord(word)
}

func (t *Trie) GetCompletion(prefix string) []string {
	values := t.head.GetPrefixWords(prefix)
	acc := make([]string, 0, len(values))

	for _, v := range values {
		acc = append(acc, prefix+v)
	}

	return acc
}

func (t *Trie) GetAllWords() []string {
	return t.head.GetAllWords()
}

func (t *Trie) Print() {
	fmt.Print("Trie:{")
	t.head.Print()
	fmt.Println("}")
}
