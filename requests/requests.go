package requests

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
)

type Result struct {
	Info  string
	Error error
}

// sendRequests sends an HTTP GET request to the specified URL.
func SendRequests(wg *sync.WaitGroup, url string, number int, results chan<- Result) {
	defer wg.Done()

	for i := 0; i < number; i++ {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Result{
				// TODO: refine this error wrapping
				Info:  "",
				Error: errors.Wrap(err, "an error occured while sending request"),
			}

			return
		}

		results <- Result{
			Info:  fmt.Sprintf("status code: %d, time: %v", resp.StatusCode, duration),
			Error: nil,
		}

	}
}
