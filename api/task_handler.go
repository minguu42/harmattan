package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minguu42/mtasks/pkg/logging"
)

type taskResponse struct {
	ID          uint64     `json:"id"`
	Title       string     `json:"title"`
	CompletedAt *time.Time `json:"completedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type tasksResponse struct {
	Tasks []*taskResponse `json:"tasks"`
}

var (
	token = "rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR"
)

type postTasksRequest struct {
	Title string `json:"title"`
}

func postTasks(w http.ResponseWriter, r *http.Request) {
	var req postTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Errorf("decoder.Decode failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(newUnauthorized(err))
		return
	}

	t, err := createTask(u.id, req.Title)
	if err != nil {
		logging.Errorf("createTask failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
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
		logging.Errorf("encoder.Encode failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newInternalServerError(err))
		return
	}
}

func getTasks(w http.ResponseWriter, _ *http.Request) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(newUnauthorized(err))
		return
	}

	ts, err := getTasksByUserID(u.id)
	if err != nil {
		logging.Errorf("getTasksByUserID failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	taskResponses := make([]*taskResponse, 0, len(ts))
	for _, t := range ts {
		tr := taskResponse{
			ID:          t.id,
			Title:       t.title,
			CompletedAt: t.completedAt,
			CreatedAt:   t.createdAt,
			UpdatedAt:   t.updatedAt,
		}
		taskResponses = append(taskResponses, &tr)
	}
	if err := encoder.Encode(tasksResponse{Tasks: taskResponses}); err != nil {
		logging.Errorf("encoder.Encode failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newInternalServerError(err))
		return
	}
}

type patchTaskRequest struct {
	IsCompleted bool `json:"isCompleted"`
}

func patchTask(w http.ResponseWriter, r *http.Request) {
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Errorf("decoder.Decode failed: %v", err)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(newUnauthorized(err))
		return
	}

	taskID, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		logging.Errorf("strconv.ParseUint failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	t, err := getTaskByID(taskID)
	if err != nil {
		logging.Errorf("getTaskByID failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	if t.userID != u.id {
		logging.Errorf("t.userID != user.id")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(newNotFound(err))
		return
	}

	resp := taskResponse{
		ID:        taskID,
		Title:     t.title,
		CreatedAt: t.createdAt,
		UpdatedAt: t.updatedAt,
	}
	if req.IsCompleted {
		now := time.Now()
		if err := updateTask(taskID, &now); err != nil {
			logging.Errorf("updateTask failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(newInternalServerError(err))
			return
		}
		resp.CompletedAt = &now
	} else {
		if err := updateTask(taskID, nil); err != nil {
			logging.Errorf("updateTask failed: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(newInternalServerError(err))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		logging.Errorf("encoder.Encode failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(newInternalServerError(err))
		return
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(newUnauthorized(err))
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		logging.Errorf("strconv.ParseUint failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	t, err := getTaskByID(id)
	if err != nil {
		logging.Errorf("getTaskByID failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(newBadRequest(err))
		return
	}

	if t.userID != u.id {
		logging.Errorf("t.userID != user.id")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(newNotFound(err))
		return
	}

	if err := destroyTask(t.id); err != nil {
		logging.Errorf("destroyTask failed: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(newInternalServerError(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
