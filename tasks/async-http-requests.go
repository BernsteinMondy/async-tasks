package tasks

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func FetchURLs(urls []string) map[string]string {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	results := make(map[string]string)

	sem := make(chan struct{}, 10)

	for _, url := range urls {
		select {
		case <-ctx.Done():
			break
		default:
			wg.Add(1)
			sem <- struct{}{}

			go func(u string) {
				defer wg.Done()
				defer func() { <-sem }()

				select {
				case <-ctx.Done():
					mu.Lock()
					results[u] = "CANCELLED: request cancelled due to another error"
					mu.Unlock()
					return
				default:
				}

				client := &http.Client{Timeout: 5 * time.Second}

				req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
				if err != nil {
					mu.Lock()
					results[u] = fmt.Sprintf("ERROR creating request: %s", err.Error())
					mu.Unlock()
					cancel()
					return
				}

				resp, err := client.Do(req)
				if err != nil {
					mu.Lock()
					results[u] = fmt.Sprintf("ERROR: %s", err.Error())
					mu.Unlock()
					cancel()
					return
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					mu.Lock()
					results[u] = fmt.Sprintf("ERROR reading body: %s", err.Error())
					mu.Unlock()
					cancel()
					return
				}

				mu.Lock()
				results[u] = fmt.Sprintf("Status: %s, Length: %d", resp.Status, len(body))
				mu.Unlock()
			}(url)
		}
	}

	wg.Wait()
	return results
}
