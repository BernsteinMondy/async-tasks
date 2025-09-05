package tasks

import (
	"sync"
	"sync/atomic"
)

func task5() {
	wg := sync.WaitGroup{}
	cnt := int64(0)

	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()
			atomic.AddInt64(&cnt, 1)
		}()
	}

	wg.Wait()
}
