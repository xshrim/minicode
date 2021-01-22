package waitdbqueue

import (
	"fmt"
	"testing"
)

type xdbNode struct {
	keyword  string
	operator byte
}

func TestQueue(t *testing.T) {

	q := NewQueue()

	q.Push(xdbNode{"sdf", 0})
	q.Push(xdbNode{"sdf", 1})
	q.Push(xdbNode{"sdf", 2})
	q.Push(xdbNode{"sdf", 3})

	q.PrintAll()

	fmt.Println("xxxxxxx:", q.PrePop())
	fmt.Println("xxxxxxx:", q.Len())
	q.PrePopClear()
	fmt.Println("xxxxxxx:", q.Len())

	n := q.Len()
	for i := 0; i < n; i++ {
		fmt.Println("============:", q.PrePop())
		//fmt.Println("长度为：",q.QueueLen())
		q.PrePopClear()
	}

	//fmt.Println("长度为：",q.QueueLen())
	//gol.Info(InsertData("sdhf1kjsdfjksdf",812381923912))

	//DeleteData("sdhf1kjsdfjksdf")
}
