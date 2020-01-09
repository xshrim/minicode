package main

import (
	"bufio"
	"flag"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
)

type Node struct {
	Char     rune
	Children []*Node
}

func NewNode(r rune) *Node {
	return &Node{Char: r}
}

func (n *Node) insert(r rune) *Node {
	child := n.get(r)
	if child == nil {
		child = NewNode(r)
		n.Children = append(n.Children, child)
	}

	return child
}

func (n *Node) get(r rune) *Node {
	for _, child := range n.Children {
		if child.Char == r {
			return child
		}
	}
	return nil
}

type Trie struct {
	Root *Node
}

func NewTrie() *Trie {
	var r rune
	trie := Trie{Root: NewNode(r)}
	return &trie
}

func (tr *Trie) Build(word string) {
	node := tr.Root
	runeArr := []rune(word)
	for _, char := range runeArr {
		child := node.insert(char)
		node = child
	}
}

func (tr *Trie) Has(word string) bool {
	node := tr.Root
	runeArr := []rune(word)
	for _, char := range runeArr {
		found := node.get(char)
		if found == nil {
			return false
		}
		node = found
	}
	return true
}

func fact(n int) int {
	if n == 0 {
		return 1
	}
	return n * fact(n-1)
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write mem profile to file")

func main() {
	flag.Parse()
	cpuProfile, _ := os.Create(*cpuprofile)
	memProfile, _ := os.Create(*memprofile)
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	var trie1 = NewTrie()
	file, err := os.Open("./test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		trie1.Build(scanner.Text())
	}

	trie1.Has("42082")
	trie1.Has("oops")
	trie1.Has("Supercalifragilisticexpialidocious")

	pprof.WriteHeapProfile(memProfile)
    
    // document
    // https://xguox.me/go-profiling-optimizing.html/
    // https://qcrao.com/2019/11/10/dive-into-go-pprof/

	// runtime and net mode

	// command
	// <go tool pprof > is equal <pprof>
	// go tool pprof <profile> ## interactive mode
	// go tool pprof -web <profile>
	// go tool pprof -http=:8888 <profile>

	// net/http/pprof
	// import _ "net/http/pprof"
	// http.ListenAndServe(":8888", nil)

	// pprof for gin
	// https://github.com/gin-contrib/pprof
}
