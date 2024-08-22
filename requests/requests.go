package requests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
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

type Work struct {
	URL               string
	N                 int
	C                 int
	DisableKeepAlives bool
}

// RunWorker sends an HTTP GET request to the specified URL.
func RunWorker(w *Work, results chan<- Result) {

	var dnsStart time.Duration

	transport := &http.Transport{
		DisableKeepAlives: w.DisableKeepAlives,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// TODO: add functions
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = now()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS lookup done: %v", now()-dnsStart)
		},
	}

	// TODO: avoid reusing requests
	// https://github.com/golang/go/issues/19653
	req, _ := http.NewRequest("GET", w.URL, nil)
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	for i := 0; i < w.N/w.C; i++ {
		start := time.Now()
		resp, err := client.Do(req)
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
