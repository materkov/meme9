package api

import (
	"context"
	"fmt"
	"strconv"
)

// /viewer
func handleViewer(_ context.Context, viewerID int, _ string) []interface{} {
	type Viewer struct {
		URL      string `json:"url,omitempty"`
		ViewerID string `json:"viewerId,omitempty"`
	}

	viewer := Viewer{
		URL: "/viewer",
	}
	results := []interface{}{&viewer}

	if viewerID != 0 {
		viewer.ViewerID = strconv.Itoa(viewerID)
		results = append(results, fmt.Sprintf("/users/%d", viewerID))
	}

	return results
}
