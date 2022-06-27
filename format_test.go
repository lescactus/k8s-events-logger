package main

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	date = "2022-06-27T14:29:05Z"
)

func TestFormatConsole(t *testing.T) {
	date, err := time.Parse(time.RFC3339, date)
	if err != nil {
		t.Fatalf(err.Error())
	}

	type args struct {
		event *corev1.Event
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Empty event",
			args: args{event: &corev1.Event{}},
			want: "Timestamp: 0001-01-01T00:00:00Z | Namespace:  | Type:  | Reason:  | Object: / | Message: ",
		},
		{
			name: "Non empty event - Created Pod",
			args: args{event: &corev1.Event{
				TypeMeta: metav1.TypeMeta{Kind: "Pod"},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unit-test",
					Namespace: "default",
				},
				FirstTimestamp: metav1.NewTime(date),
				Type:           "Normal",
				Reason:         "Created",
				Message:        "Created container unit-test",
			}},
			want: "Timestamp: 2022-06-27T14:29:05Z | Namespace: default | Type: Normal | Reason: Created | Object: / | Message: Created container unit-test",
		},
		{
			name: "Non empty event - Killing Pod",
			args: args{event: &corev1.Event{
				TypeMeta: metav1.TypeMeta{Kind: "Pod"},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unit-test",
					Namespace: "default",
				},
				FirstTimestamp: metav1.NewTime(date),
				Type:           "Normal",
				Reason:         "Killing",
				Message:        "Killing container unit-test",
			}},
			want: "Timestamp: 2022-06-27T14:29:05Z | Namespace: default | Type: Normal | Reason: Killing | Object: / | Message: Killing container unit-test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatConsole(tt.args.event); got != tt.want {
				fmt.Println(got)
				t.Errorf("FormatConsole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatJSON(t *testing.T) {
	date, err := time.Parse(time.RFC3339, date)
	if err != nil {
		t.Fatalf(err.Error())
	}

	type args struct {
		event *corev1.Event
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Empty event",
			args:    args{event: &corev1.Event{}},
			want:    []byte(`{"metadata":{"creationTimestamp":null},"involvedObject":{},"source":{},"firstTimestamp":null,"lastTimestamp":null,"eventTime":null,"reportingComponent":"","reportingInstance":""}`),
			wantErr: false,
		},
		{
			name: "Non empty event - Created Pod",
			args: args{event: &corev1.Event{
				TypeMeta: metav1.TypeMeta{Kind: "Pod"},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "unit-test",
					Namespace: "default",
				},
				FirstTimestamp: metav1.NewTime(date),
				Type:           "Normal",
				Reason:         "Created",
				Message:        "Created container unit-test",
			}},
			want:    []byte(`{"kind":"Pod","metadata":{"name":"unit-test","namespace":"default","creationTimestamp":null},"involvedObject":{},"reason":"Created","message":"Created container unit-test","source":{},"firstTimestamp":"2022-06-27T14:29:05Z","lastTimestamp":null,"type":"Normal","eventTime":null,"reportingComponent":"","reportingInstance":""}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatJSON(tt.args.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("FormatJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FormatJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
