package tasks

import "sync"

func task3() <-chan int {
	ch := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := range 5 {
			ch <- i
		}
	}()

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}
