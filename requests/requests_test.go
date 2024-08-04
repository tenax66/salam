package requests

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSendRequests(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Test!"))
	}))
	defer ts.Close()

	type args struct {
		wg      *sync.WaitGroup
		url     string
		number  int
		results chan Result
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "normal",
			args: args{
				&sync.WaitGroup{},
				ts.URL,
				5,
				make(chan Result, 5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.wg.Add(1)
			SendRequests(tt.args.wg, tt.args.url, tt.args.number, tt.args.results)

			tt.args.wg.Wait()
			close(tt.args.results)

			for result := range tt.args.results {
				if (result.Error) != nil {
					t.Errorf("expected no error, got %v", result.Error)
				}
			}
		})
	}
}
