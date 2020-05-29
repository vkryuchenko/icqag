package webhook

import (
	"net/http"
)

// Payload interface for any data from any alert systems
type Payload interface {
	Parse(req *http.Request) (string, error)
}
