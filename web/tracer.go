package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Tracer struct {
	Started time.Time
	Name    string
	ID      int
	TraceID int
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
	WriteSpan(t.TraceID, t.Name, t.Started)
}

func (t *Tracer) StartChild(name string) *Tracer {
	return &Tracer{
		Started: time.Now(),
		Name:    name,
		TraceID: t.TraceID,
		ID:      rand.Int(),
	}
}

func WriteSpan(traceID int, name string, started time.Time) {
	const template = `
		[
			{
				"id": "%x",
				"traceId": "%x",
				"name": "%s",
				"timestamp": %d,
				"duration": %d,
				"localEndpoint": {
					"serviceName": "web"
				}
			}
		]
	`
	body := fmt.Sprintf(template, rand.Int(), traceID, name, started.UnixNano()/1000, time.Since(started).Microseconds())

	go func() {
		_, _ = http.Post("http://127.0.0.1:9411/api/v2/spans", "application/json", strings.NewReader(body))
	}()
}
