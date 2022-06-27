package main

import (
	"encoding/json"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
)

func FormatConsole(event *corev1.Event) string {
	return fmt.Sprintf("Timestamp: %v | Namespace: %s | Type: %s | Reason: %s | Object: %s/%s | Message: %s",
		event.FirstTimestamp.Format(time.RFC3339),
		event.Namespace,
		event.Type,
		event.Reason,
		event.InvolvedObject.Kind,
		event.InvolvedObject.Name,
		event.Message,
	)
}

func FormatJSON(event *corev1.Event) ([]byte, error) {
	e, err := json.Marshal(event)
	if err != nil {
		return nil, fmt.Errorf("failed to format event in json: %w", err)
	}

	return e, nil
}
