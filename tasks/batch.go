package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	maxBatchSize = 5
	batchTimeout = 2 * time.Second
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		StartBatchProcessor(ctx, input)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(input)

		for i := 1; i <= 20; i++ {
			select {
			case <-ctx.Done():
				fmt.Println("Sender: context canceled, stopping")
				return
			case input <- i:
				time.Sleep(300 * time.Millisecond)
			}
		}
		fmt.Println("Sender: all data sent")
	}()

	<-ctx.Done()
	fmt.Println("Main: context deadline exceeded, waiting for graceful shutdown...")

	wg.Wait()
	fmt.Println("Main: all goroutines completed gracefully")
}

func StartBatchProcessor(ctx context.Context, input <-chan int) {
	timer := time.NewTimer(batchTimeout)
	defer timer.Stop()

	batch := make([]int, 0, maxBatchSize)
	defer processRemainingBatch(batch)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("BatchProcessor: context canceled, exiting")
			return

		case <-timer.C:
			if len(batch) > 0 {
				fmt.Println("BatchProcessor: timed out")
				fmt.Println("Processed batch:", batch)
				batch = batch[:0]
			}
			timer.Reset(batchTimeout)

		case val, ok := <-input:
			if !ok {
				fmt.Println("BatchProcessor: input channel closed, exiting")
				return
			}

			batch = append(batch, val)
			if len(batch) == maxBatchSize {
				fmt.Println("Processed batch:", batch)
				batch = batch[:0]

				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(batchTimeout)
			}
		}
	}
}

func processRemainingBatch(batch []int) {
	if len(batch) > 0 {
		fmt.Println("Processing remaining batch:", batch)
	}
}
