package http

import "net/http"

func (h *Handler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeJSONData(w, http.StatusOK, "API is UP!")
}
