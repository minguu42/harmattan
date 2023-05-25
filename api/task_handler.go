package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type taskResponse struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type postTasksRequest struct {
	Title string `json:"title"`
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	var req postTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Errorf("decoder.Decode failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
	}

	token := "rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR"
	u, err := getUserByToken(token)
	if err != nil {
		Errorf("getUserByToken failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
	}

	t, err := createTask(u.id, req.Title)
	if err != nil {
		Errorf("createTask failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("http://localhost:8080/tasks/%d", t.id))
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	resp := taskResponse{
		ID:          t.id,
		Title:       t.title,
		CompletedAt: t.completedAt,
		CreatedAt:   t.createdAt,
		UpdatedAt:   t.updatedAt,
	}
	if err := encoder.Encode(resp); err != nil {
		Errorf("encoder.Encode failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newInternalServerError(err))
	}
}
