package b

import (
	sub2 "git.oschina.net/jscode/list-interfaces/testdata/b/sub"
	a "sync"
	"git.oschina.net/jscode/list-interfaces/testdata/b/suba"
)

type SB struct {
}

func (this  SB) Add(a sub2.SubSA, locker a.Locker, b B, subsa1 suba.SubSa1){}



