package api

const (
	TypeRequest = iota
	TypeResult
)

type dispatcherMsg struct {
	Type int
	URL  string
	Data interface{}
}

func DoAsync(requestID, viewerID int, urls []string) []interface{} {
	done := make(chan []interface{})
	dispatcherChan := make(chan dispatcherMsg)

	go func() {
		results := map[string]interface{}{}
		waiting := 0

		for {
			op := <-dispatcherChan

			if op.Type == TypeRequest {
				if _, isResolved := results[op.URL]; isResolved {
					continue
				}

				waiting++
				go func() {
					resource, related := handleResource(requestID, viewerID, op.URL)

					for _, url := range related {
						dispatcherChan <- dispatcherMsg{
							Type: TypeRequest,
							URL:  url,
						}
					}

					dispatcherChan <- dispatcherMsg{
						Type: TypeResult,
						URL:  op.URL,
						Data: resource,
					}
				}()
			} else if op.Type == TypeResult {
				results[op.URL] = op.Data
				waiting--

				if waiting == 0 {
					resultsList := make([]interface{}, len(results))

					idx := 0
					for _, item := range results {
						resultsList[idx] = item
						idx++
					}

					done <- resultsList
					break
				}
			}
		}
	}()

	for _, url := range urls {
		dispatcherChan <- dispatcherMsg{
			Type: TypeRequest,
			URL:  url,
		}
	}

	return <-done
}
