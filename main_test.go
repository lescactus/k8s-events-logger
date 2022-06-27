package main

import (
	"os"
	"testing"
)

func TestIsInNamespacesToWatch(t *testing.T) {
	type args struct {
		namespace         string
		namespacesToWatch []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "One namespace to watch - given namespace is to watch",
			args: args{
				namespace:         "default",
				namespacesToWatch: []string{"default"},
			},
			want: true,
		},
		{
			name: "Several namespaces to watch - given namespace is to watch",
			args: args{
				namespace:         "default",
				namespacesToWatch: []string{"default", "ns1", "ns2", "ns3"},
			},
			want: true,
		},
		{
			name: "One namespace to watch - given namespace is not to watch",
			args: args{
				namespace:         "default",
				namespacesToWatch: []string{"ns1"},
			},
			want: false,
		},
		{
			name: "Several namespaces to watch - given namespace is not to watch",
			args: args{
				namespace:         "default",
				namespacesToWatch: []string{"ns1", "ns2", "ns3"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInNamespacesToWatch(tt.args.namespace, tt.args.namespacesToWatch); got != tt.want {
				t.Errorf("IsInNamespacesToWatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRestConfig(t *testing.T) {
	validKubeConfig := "testsdata/kubeconfig/config.valid"
	invalidKubeConfig := "testsdata/kubeconfig/config.invalid"

	t.Run("Valid kubeconfig -"+validKubeConfig, func(t *testing.T) {
		os.Setenv("KUBECONFIG", validKubeConfig)

		_, err := NewRestConfig()
		if err != nil {
			t.Errorf("NewRestConfig() error = %v", err)
			return
		}
	})

	t.Run("Invalid kubeconfig -"+invalidKubeConfig, func(t *testing.T) {
		os.Setenv("KUBECONFIG", invalidKubeConfig)

		_, err := NewRestConfig()
		if err == nil {
			t.Errorf("NewRestConfig() problem: Expected error, got nil")
			return
		}
	})
}
