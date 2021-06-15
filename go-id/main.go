package main

import (
	"fmt"
	"sync"

	"github.com/go-id/id"
)

var wg sync.WaitGroup

func main() {
	w := id.NewWorker(1, 1)

	ch := make(chan uint64, 1000)
	defer close(ch)
	count := 5
	wg.Add(count)

	//并发ID生成
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id, _ := w.NextID()
			ch <- id
		}()
	}
	wg.Wait()
	m := make(map[uint64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		// 检查ID重复
		_, err := m[id]
		if err {
			fmt.Printf("repeat id %d\n", id)
			return
		}
		// id存入map
		m[id] = i
	}
	fmt.Println("All", len(m), " ID Get!")
	//fmt.Println(m)

}
