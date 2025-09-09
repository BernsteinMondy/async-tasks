package tasks

import (
	"fmt"
	"sync"
)

func task2() {
	wg := &sync.WaitGroup{}

	wg.Add(5)
	for i := range 5 {
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}

	wg.Wait()
}

// 1. Канал с буфером, ждем пока заполнится буфер.
// 2. Mutex с счетчиком
