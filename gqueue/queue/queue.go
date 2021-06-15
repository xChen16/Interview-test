package queue

import (
	"sync"
	"time"
)

var (
	firstSliceSize       = 1   //第一个切片Size
	maxFirstSliceSize    = 16  //切片最大Size
	maxInternalSliceSize = 128 //内存切片最大Size
)

const (
	signalInterval = 200
	signalChanSize = 10
)

// FIFO队列
type Queue struct {
	head          *Node //	头指针
	tail          *Node //	尾指针
	hp            int   //	第一个值
	len           int   //	队列长度
	lastSliceSize int   //	last切片Size
	sync.Mutex

	C chan struct{}
}

// 队列节点
type Node struct {
	v []interface{}
	n *Node
}

// 初始化队列
func New() *Queue {
	q := new(Queue).Init()
	//fmt.Printf("%d", q.len)
	go func() {
		ticker := time.NewTicker(time.Millisecond * signalInterval)
		defer ticker.Stop()
		select {
		case <-ticker.C:
			if q.Len() > 0 {
				select {
				case q.C <- struct{}{}:
				default:
				}
			}
		}
	}()
	return q
}

func (q *Queue) Init() *Queue {
	q.head = nil
	q.tail = nil
	q.hp = 0
	q.len = 0
	return q
}

// 队列长度
func (q *Queue) Len() int {
	q.Lock()
	l := q.len
	q.Unlock()
	return l
}

// 队列判空
func (q *Queue) Front() (interface{}, bool) {
	q.Lock()
	defer q.Unlock()
	if q.head == nil {
		return nil, false
	}
	return q.head.v[q.hp], true
}

// 初始化节点
func newNode(capacity int) *Node {
	return &Node{
		v: make([]interface{}, 0, capacity),
	}
}

// 队列push
func (q *Queue) Push(v interface{}) {
	q.Lock()
	defer q.Unlock()
	if q.head == nil {
		h := newNode(firstSliceSize)
		q.head = h
		q.tail = h
		q.lastSliceSize = maxFirstSliceSize
	} else if len(q.tail.v) >= q.lastSliceSize {
		n := newNode(maxInternalSliceSize)
		q.tail.n = n
		q.tail = n
		q.lastSliceSize = maxInternalSliceSize
	}
	q.tail.v = append(q.tail.v, v)
	q.len++
}

// 队列pop
func (q *Queue) Pop() (interface{}, bool) {
	q.Lock()
	defer q.Unlock()
	if q.head == nil {
		return nil, false
	}
	v := q.head.v[q.hp]
	q.head.v[q.hp] = nil
	q.len--
	q.hp++
	if q.hp >= len(q.head.v) {
		n := q.head.n
		q.head.n = nil
		q.head = n
		q.hp = 0
	}
	return v, true
}
