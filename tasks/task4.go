package tasks

import (
	"fmt"
	"sync"
)

func task4() {
	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}
	cnt := 0

	wg.Add(10)
	for range 10 {
		go func() {
			defer wg.Done()

			m.Lock()
			cnt++
			m.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(cnt)
}
