package tasks

import "sync"

func Merge(cs ...<-chan int) <-chan int {
	wg := new(sync.WaitGroup)
	out := make(chan int)

	fn := func(c <-chan int) {
		defer wg.Done()

		for n := range c {
			out <- n
		}
	}

	for _, c := range cs {
		wg.Add(1)
		go fn(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
