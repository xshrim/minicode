package main

import (
	"git.oschina.net/jscode/list-interfaces/testdata/b"
	"git.oschina.net/jscode/list-interfaces/testdata/b/sub"
)

func m1(a1 b.IA) {
	a1.Add(sub2.SubSA{})
}

func main() {
	m1(&b.SB{})
}

