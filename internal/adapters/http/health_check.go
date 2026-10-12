package http

import "net/http"

func (h *Handler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.healthcheck.Health(r.Context()); err != nil {
		_ = writeJSONError(w, http.StatusServiceUnavailable, "service unavailable")
		return
	}
	_ = writeJSONData(w, http.StatusOK, "API is up")
}
