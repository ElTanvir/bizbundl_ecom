package admin

import (
	"testing"
)

// go test ./internal/views/frontend/...
func TestRequiredRoutesExist(t *testing.T) {
	requiredPaths := []string{
		"/",
	}

	for _, required := range requiredPaths {
		found := false
		for _, route := range Routes {
			if route.Path == required {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Critical route missing: %s. Did you accidentally delete it from Routes?", required)
		}
	}
}
