package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/minguu42/mtasks/pkg/app/render"
	"github.com/minguu42/mtasks/pkg/logging"
)

var (
	token = "rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR"
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

type postTasksRequest struct {
	Title string `json:"title"`
}

// PostTasks は POST /tasks に対応するハンドラ関数
func PostTasks(w http.ResponseWriter, r *http.Request) {
	var req postTasksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Errorf("decoder.Decode failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		render.Error(w, http.StatusUnauthorized, err)
		return
	}

	t, err := createTask(u.id, req.Title)
	if err != nil {
		logging.Errorf("createTask failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("http://localhost:8080/tasks/%d", t.id))
	resp := taskResponse{
		ID:          t.id,
		Title:       t.title,
		CompletedAt: t.completedAt,
		CreatedAt:   t.createdAt,
		UpdatedAt:   t.updatedAt,
	}
	if err := render.Response(w, http.StatusCreated, resp); err != nil {
		logging.Errorf("encoder.Encode failed: %v", err)
		render.Error(w, http.StatusInternalServerError, err)
		return
	}
}

// GetTasks は GET /tasks に対応するハンドラ関数
func GetTasks(w http.ResponseWriter, _ *http.Request) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		render.Error(w, http.StatusUnauthorized, err)
		return
	}

	ts, err := getTasksByUserID(u.id)
	if err != nil {
		logging.Errorf("getTasksByUserID failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

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
	if err := render.Response(w, http.StatusOK, tasksResponse{Tasks: taskResponses}); err != nil {
		logging.Errorf("encoder.Encode failed: %v", err)
		render.Error(w, http.StatusInternalServerError, err)
		return
	}
}

type patchTaskRequest struct {
	IsCompleted bool `json:"isCompleted"`
}

// PatchTask は PATCH /tasks/{taskID} に対応するハンドラ関数
func PatchTask(w http.ResponseWriter, r *http.Request) {
	var req patchTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logging.Errorf("decoder.Decode failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		render.Error(w, http.StatusUnauthorized, err)
		return
	}

	taskID, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		logging.Errorf("strconv.ParseUint failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	t, err := getTaskByID(taskID)
	if err != nil {
		logging.Errorf("getTaskByID failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	if t.userID != u.id {
		logging.Errorf("t.userID != user.id")
		render.Error(w, http.StatusNotFound, err)
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
			render.Error(w, http.StatusInternalServerError, err)
			return
		}
		resp.CompletedAt = &now
	} else {
		if err := updateTask(taskID, nil); err != nil {
			logging.Errorf("updateTask failed: %v", err)
			render.Error(w, http.StatusInternalServerError, err)
			return
		}
	}

	if err = render.Response(w, http.StatusOK, resp); err != nil {
		logging.Errorf("encoder.Encode failed: %v", err)
		render.Error(w, http.StatusInternalServerError, err)
		return
	}
}

// DeleteTask は DELETE /tasks/{taskID} に対応するハンドラ関数
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	u, err := getUserByToken(token)
	if err != nil {
		logging.Errorf("getUserByToken failed: %v", err)
		render.Error(w, http.StatusUnauthorized, err)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "taskID"), 10, 64)
	if err != nil {
		logging.Errorf("strconv.ParseUint failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	t, err := getTaskByID(id)
	if err != nil {
		logging.Errorf("getTaskByID failed: %v", err)
		render.Error(w, http.StatusBadRequest, err)
		return
	}

	if t.userID != u.id {
		logging.Errorf("t.userID != user.id")
		render.Error(w, http.StatusNotFound, err)
		return
	}

	if err := destroyTask(t.id); err != nil {
		logging.Errorf("destroyTask failed: %v", err)
		render.Error(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
