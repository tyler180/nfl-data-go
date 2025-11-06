package download

import "fmt"

// HTTPError represents a non-2xx HTTP status with a small body preview.
type HTTPError struct {
	Code int
	Body string
}

func (e *HTTPError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Body == "" {
		return fmt.Sprintf("http error: %d", e.Code)
	}
	return fmt.Sprintf("http error: %d: %s", e.Code, e.Body)
}
