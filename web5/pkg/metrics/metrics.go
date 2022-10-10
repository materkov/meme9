package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

type span struct {
	Id            string `json:"id"`
	TraceId       string `json:"traceId"`
	Timestamp     int64  `json:"timestamp"`
	Duration      int64  `json:"duration"`
	Name          string `json:"name"`
	LocalEndpoint struct {
		ServiceName string `json:"serviceName"`
	} `json:"localEndpoint"`
}

func WriteSpan(method string, traceID int, main bool, duration time.Duration) error {
	id := fmt.Sprintf("%016x", rand.Int63())
	if main {
		id = fmt.Sprintf("%016x", traceID)
	}
	spans := []span{
		{
			Id:        id,
			TraceId:   fmt.Sprintf("%016x", traceID),
			Timestamp: time.Now().Add(-duration).UnixMicro(),
			Duration:  duration.Microseconds(),
			Name:      method,
		},
	}

	spans[0].LocalEndpoint.ServiceName = "web"

	spansBytes, _ := json.Marshal(spans)

	resp, err := http.Post("http://localhost:9411/api/v2/spans", "application/json", bytes.NewReader(spansBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 202 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("incorrect http response: %d, %s", resp.StatusCode, body)
	}

	return nil
}
