package main

import (
	"fmt"
	"sync"

	"github.com/gqueue/queue"
)

func main() {
	var q = queue.New()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < 10; i++ {
			q.Push(i + 1)

		}
		wg.Done()
	}()
	// 消费者.
	go func() {
	LOOP:
		for {
			select {
			case <-q.C:
				for {
					i, ok := q.Pop()
					if !ok {
						// no msg available
						continue LOOP
					}

					fmt.Printf("%d\n", i.(int))
				}
			}

		}

	}()

	wg.Wait()
}
