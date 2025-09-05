package tasks

import (
	"fmt"
	"sync"
)

func task1() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Hello from goroutine!")
	}()

	wg.Wait()
}

// 1. Done channel - ждать сигнала оттуда
// 2. Mutex с флагом и бесконечный for
