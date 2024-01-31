package tracer

import (
	"bytes"
	"context"
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
	TraceID int
	Tags    map[string]string
}

func NewTracer(name string) *Tracer {
	return &Tracer{
		Started: time.Now(),
		Name:    name,
		TraceID: rand.Int(),
		Tags:    map[string]string{},
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
		Tags:    map[string]string{},
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

	//log.Printf("Send: %s", spansBytes)

	go func() {
		_, _ = http.Post("http://127.0.0.1:9411/api/v2/spans", "application/json", bytes.NewReader(spansBytes))
	}()
}

func FromCtx(ctx context.Context) *Tracer {
	tracer, ok := ctx.Value(ctxKey).(*Tracer)
	if !ok {
		tracer = &Tracer{}
	}
	return tracer
}

func WithCtx(ctx context.Context, t *Tracer) context.Context {
	return context.WithValue(ctx, ctxKey, t)
}

type contextKey string

var ctxKey contextKey = "tracer"
