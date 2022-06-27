package main

import (
	"testing"

	"github.com/spf13/viper"
)

func TestAppConfigIsValidOutput(t *testing.T) {
	type fields struct {
		Viper *viper.Viper
	}
	tests := []struct {
		name   string
		fields fields
		output string
		want   bool
	}{
		{
			name:   "Valid OUTPUT value - console",
			fields: fields{Viper: viper.New()},
			output: "console",
			want:   true,
		},
		{
			name:   "Valid OUTPUT value - json",
			fields: fields{Viper: viper.New()},
			output: "json",
			want:   true,
		},
		{
			name:   "Invalid OUTPUT value - invalid",
			fields: fields{Viper: viper.New()},
			output: "invalid",
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &AppConfig{
				Viper: tt.fields.Viper,
			}
			config.SetDefault("OUTPUT", tt.output)

			if got := config.isValidOutput(); got != tt.want {
				t.Errorf("AppConfig.isValidOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}
