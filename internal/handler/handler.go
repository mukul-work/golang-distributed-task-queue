package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mukul-work/golang-distributed-task-queue/authutil"
	"github.com/mukul-work/golang-distributed-task-queue/internal/dbgen"
	"github.com/mukul-work/golang-distributed-task-queue/internal/service"
	"github.com/mukul-work/golang-distributed-task-queue/models"
)

type Handler struct {
	q service.Queue
}

func NewHandlerService(q service.Queue) *Handler {
	if q == nil {
		return nil
	}
	return &Handler{
		q: q,
	}
}

func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) PostJob(w http.ResponseWriter, r *http.Request) {

	var req models.JobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, models.ErrorMessage{
			Message: "Bad Request",
			Error:   err.Error(),
		})
		return
	}

	if len(req.Payload) == 0 {
		writeJson(w, http.StatusBadRequest, models.ErrorMessage{
			Message: "Payload should not be empty",
		})
		return
	}

	id := uuid.New().String()
	job := dbgen.Job{
		ID:          id,
		Type:        req.Type,
		Payload:     req.Payload,
		Status:      "pending",
		MaxAttempts: 4,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := h.q.Enqueue(ctx, job); err != nil {
		writeJson(w, http.StatusInternalServerError, models.ErrorMessage{
			Message: "Failed to register the job",
			Error:   err.Error(),
		})
		return
	}

	writeJson(w, http.StatusOK, models.SuccessMessage{
		Message: "Job registered successfully",
		ID:      fmt.Sprintf("ID: %s", job.ID),
	})
}

func (h *Handler) PostAdminApiKey(w http.ResponseWriter, r *http.Request) {

	var req models.PostAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, models.ErrorMessage{
			Message: "Invalid Request",
			Error:   err.Error(),
		})
		return
	}

	key, err := authutil.KeyGenerator(req.Prefix)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, models.ErrorMessage{
			Message: "Failed to generate the API key",
		})
		return
	}
	hashKey := authutil.HashApiKeys(key)

	newKey, err := h.q.CreateAPIKey(r.Context(), dbgen.CreateAPIKeyParams{
		ID:      uuid.New().String(),
		Name:    req.Name,
		KeyHash: hashKey,
	})

	if err != nil {
		writeJson(w, http.StatusInternalServerError, models.ErrorMessage{
			Message: "Failedto generate API key",
		})
		return
	}

	writeJson(w, http.StatusCreated, map[string]any{
		"id":         newKey.ID,
		"name":       newKey.Name,
		"key":        key,
		"created_at": newKey.CreatedAt,
		"warning":    "Save this key securely. It will not be visible again",
	})
}

func (h *Handler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	job, err := h.q.GetJob(r.Context(), id)
	if err != nil {
		writeJson(w, http.StatusNotFound, models.ErrorMessage{
			Message: "ID could not be found",
		})
	}
	writeJson(w, http.StatusOK, job)
}

func (h *Handler) GetApiKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := h.q.ListAPIKeys(r.Context())
	if err != nil {
		writeJson(w, http.StatusInternalServerError, models.ErrorMessage{
			Message: "Could not List the API keys",
			Error:   err.Error(),
		})
		return
	}

	writeJson(w, http.StatusOK, map[string]any{
		"Message": "These are the keys",
		"Keys":    keys,
	})
}
