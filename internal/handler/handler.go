package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mukul-work/golang-distributed-task-queue/models"
)

type Handler struct {
	Queue chan models.Job
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var job models.Job

	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Failed to parse data", http.StatusBadRequest)
		return
	}

	select {
	case h.Queue <- job:
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Job Accepted"))
	default:
		http.Error(w, "Queue is full", http.StatusServiceUnavailable)
		return

	}

}
