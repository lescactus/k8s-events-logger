package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHealthzServer(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
		want *HealthzServer
	}{
		{
			name: "Address - 127.0.0.1:80",
			args: args{addr: "127.0.0.1:80"},
			want: &HealthzServer{server: http.Server{Addr: "127.0.0.1:80", ReadTimeout: 5 * time.Second, ReadHeaderTimeout: 3 * time.Second, WriteTimeout: 5 * time.Second}},
		},
		{
			name: "Address - :8080",
			args: args{addr: "8080"},
			want: &HealthzServer{server: http.Server{Addr: "8080", ReadTimeout: 5 * time.Second, ReadHeaderTimeout: 3 * time.Second, WriteTimeout: 5 * time.Second}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHealthzServer(tt.args.addr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHealthzServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthzServerHealthz(t *testing.T) {
	s := NewHealthzServer(":8080")
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/ready", "/alive":
			s.Healthz(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
	defer mockServer.Close()

	resp, err := http.Get(fmt.Sprintf("%s/ready", mockServer.URL))
	if err != nil {
		t.Fatal(err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":200,"message":"Server ready and alive"}`), data)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
