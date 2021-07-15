package tracer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Tracer struct {
	Started time.Time
	Name    string
	ID      int
	TraceID int
	Tags    map[string]string
}

func NewTracer(name string) *Tracer {
	return &Tracer{
		Started: time.Now(),
		Name:    name,
		TraceID: rand.Int(),
		ID:      rand.Int(),
	}
}

func (t *Tracer) Stop() {
	WriteSpan(t.TraceID, t.Name, t.Started, t.Tags)
}

func (t *Tracer) StartChild(name string) *Tracer {
	return &Tracer{
		Started: time.Now(),
		Name:    name,
		TraceID: t.TraceID,
		ID:      rand.Int(),
	}
}

func WriteSpan(traceID int, name string, started time.Time, tags map[string]string) {
	type LocalEndpoint struct {
		ServiceName string `json:"serviceName"`
	}

	type Span struct {
		ID            string            `json:"id"`
		TraceID       string            `json:"traceId"`
		Name          string            `json:"name"`
		Timestamp     int64             `json:"timestamp"`
		Duration      int               `json:"duration"`
		LocalEndpoint LocalEndpoint     `json:"localEndpoint"`
		Tags          map[string]string `json:"tags,omitempty"`
	}

	spans := []Span{{
		ID:            fmt.Sprintf("%x", rand.Int()),
		TraceID:       fmt.Sprintf("%x", traceID),
		Name:          name,
		Timestamp:     started.UnixNano() / 1000,
		Duration:      int(time.Since(started).Microseconds()),
		LocalEndpoint: LocalEndpoint{ServiceName: "web"},
		Tags:          tags,
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
