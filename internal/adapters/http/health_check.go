package http

import (
	"fmt"
	"net/http"
)

func (h *Handler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.healthcheck.Health(r.Context()); err != nil {
		_ = writeJSONError(w, http.StatusServiceUnavailable, fmt.Sprintf("service unavaliable: %v", err))
		return
	}
	_ = writeJSONData(w, http.StatusOK, "API and all its services are up")
}
