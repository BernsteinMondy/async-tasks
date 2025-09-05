package tasks

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := make(chan int)

	go StartBatchProcessor(ctx, input)

	go func() {
		for i := 1; i <= 20; i++ {
			input <- i
			time.Sleep(300 * time.Millisecond)
		}
	}()

	<-ctx.Done()
	fmt.Println("Main: processing stopped")
}

func StartBatchProcessor(ctx context.Context, input <-chan int) {
	const (
		maxBatchSize = 5
		batchTimeout = 2 * time.Second
	)

	timer := time.NewTimer(batchTimeout)
	batch := make([]int, 0, maxBatchSize)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("BatchProcessor: context canceled")
		case <-timer.C:
			if len(batch) > 0 {
				fmt.Println("BatchProcessor: timed out")
				fmt.Println("Processed batch:", batch)
				batch = batch[:0]
			}

			timer.Reset(batchTimeout)
		case val := <-input:
			batch = append(batch, val)
			if len(batch) == maxBatchSize {
				fmt.Println("Processed batch:", batch)
				batch = batch[:0]
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(batchTimeout)
			}
		}
	}
}
