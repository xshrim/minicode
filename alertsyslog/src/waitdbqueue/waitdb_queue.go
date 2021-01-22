package waitdbqueue

import (
	"container/list"
	"sync"

	"github.com/xshrim/gol"
)

var lock sync.Mutex

// const N int = 10

type QueueList struct {
	/*	figure  int
		digits1 [N]int
		digits2 [N]int
		sflag   bool*/

	data *list.List
}

// 新建队列
func NewQueue() *QueueList {
	q := new(QueueList)
	q.data = list.New()
	return q
}

// 向队列尾部推送数据
func (q *QueueList) Push(v interface{}) {
	defer lock.Unlock()
	lock.Lock()
	q.data.PushFront(v)
}

// 遍历队列信息，但不清理
func (q *QueueList) PrintAll() {
	gol.Info("==========数据库中断期间，内存信息保留如下：===========")
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		// fmt.Println("item:", iter.Value)
		gol.Info(iter.Value)
	}
	gol.Info("================================================")
}

// 获取队列头部数据，并清理该数据
func (q *QueueList) Pop() interface{} {
	defer lock.Unlock()
	lock.Lock()
	iter := q.data.Back()
	v := iter.Value
	q.data.Remove(iter)
	return v
}

// 获取队列头部信息，暂不清理
func (q *QueueList) PrePop() interface{} {
	iter := q.data.Back()
	v := iter.Value
	//q.data.Remove(iter)
	return v
}

// 清理队列头部信息
func (q *QueueList) PrePopClear() {
	defer lock.Unlock()
	lock.Lock()
	iter := q.data.Back()
	q.data.Remove(iter)
}

func (q *QueueList) Len() int {
	return q.data.Len()
}
