package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func WriteSpan(requestID int, name string, duration time.Duration, tags ...string) {
	type LocalEndpoint struct {
		ServiceName string `json:"serviceName"`
	}
	type Span struct {
		Debug         bool              `json:"debug"`
		ID            string            `json:"id"`
		TraceID       string            `json:"traceId"`
		Name          string            `json:"name"`
		Timestamp     int64             `json:"timestamp"`
		Duration      int64             `json:"duration"`
		LocalEndpoint LocalEndpoint     `json:"localEndpoint"`
		Tags          map[string]string `json:"tags"`
	}

	tagsMap := map[string]string{}
	for i := 0; i < len(tags); i += 2 {
		tagsMap[tags[i]] = tags[i+1]
	}

	spans := []Span{{
		Debug:         true,
		ID:            fmt.Sprintf("%x", rand.Int()),
		TraceID:       fmt.Sprintf("%x", requestID),
		Name:          name,
		Timestamp:     time.Now().Add(-duration).UnixMicro(),
		Duration:      duration.Microseconds(),
		LocalEndpoint: LocalEndpoint{ServiceName: "web"},
		Tags:          tagsMap,
	}}

	spansBytes, err := json.Marshal(spans)
	if err != nil {
		log.Printf("Error json: %s", err)
		return
	}

	log.Printf("Send: %s", spansBytes)

	go func() {
		_, _ = http.Post("http://127.0.0.1:9411/api/v2/spans", "application/json", bytes.NewReader(spansBytes))
	}()
}