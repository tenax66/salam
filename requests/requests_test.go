package requests

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendRequests(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, Test!"))
	}))
	defer ts.Close()

	type args struct {
		url     string
		number  int
		results chan Result
	}
	tests := []struct {
		name      string
		args      args
		expectErr bool
	}{
		{
			name: "normal",
			args: args{
				ts.URL,
				5,
				make(chan Result, 5),
			},
			expectErr: false,
		},
		{
			name: "invalid url",
			args: args{
				"abc://xyz",
				5,
				make(chan Result, 5),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendRequests(tt.args.url, tt.args.number, tt.args.results)
			close(tt.args.results)

			for result := range tt.args.results {
				if tt.expectErr == false && result.Error != nil {
					t.Errorf("expected no error, got %v", result.Error)
				} else if tt.expectErr == true && result.Error == nil {
					t.Errorf("expected error, got nil")
				}
			}
		})
	}
}
