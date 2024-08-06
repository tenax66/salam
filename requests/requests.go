package requests

import (
	"io"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
)

// Result represents the result of a HTTP request.
type Result struct {
	StatusCode int
	Body       string
	Duration   time.Duration
	Error      error
}

// SendRequests sends an HTTP GET request to the specified URL.
func SendRequests(url string, number int, results chan<- Result) {
	for i := 0; i < number; i++ {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)

		if err != nil {
			results <- Result{
				Error: errors.Wrap(err, "error while sending request"),
			}

			return
		}

		// Prevent resource leaks and enable keep-alive
		// cf. https://pkg.go.dev/net/http#Client
		// > If the Body is not both read to EOF and closed,
		// > the Client's underlying RoundTripper (typically Transport) may not be able to
		// > re-use a persistent TCP connection to the server for a subsequent "keep-alive" request.

		defer resp.Body.Close()
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			results <- Result{
				Error: errors.Wrap(err, "error while reading the response body"),
			}

			return
		}

		results <- Result{
			StatusCode: resp.StatusCode,
			Body:       string(b),
			Duration:   duration,
			Error:      nil,
		}

	}
}
