package app

import (
	"context"

	"github.com/minguu42/mtasks/app/ogen"
)

var token = "rAM9Fm9huuWEKLdCwHBcju9Ty_-TL2tDsAicmMrXmUnaCGp3RtywzYpMDPdEtYtR"

// PostTasks は POST /projects/{projectID}/tasks に対応するハンドラ
func (h *handler) PostTasks(_ context.Context, _ *ogen.PostTasksReq, _ ogen.PostTasksParams) (ogen.PostTasksRes, error) {
	return &ogen.PostTasksNotImplemented{}, nil
	//u, err := h.repository.GetUserByToken(ctx, token)
	//if err != nil {
	//	logging.Errorf("getUserByToken failed: %v", err)
	//	return &ogen.PostTasksUnauthorized{}, nil
	//}
	//
	//t, err := h.repository.CreateTask(ctx, u.ID, req.Title)
	//if err != nil {
	//	logging.Errorf("createTask failed: %v", err)
	//	return &ogen.PostTasksBadRequest{}, nil
	//}
	//
	//return &ogen.TaskHeaders{
	//	Location: fmt.Sprintf("http://localhost:8080/tasks/%d", t.ID),
	//	Response: newTaskResponse(t),
	//}, nil
}

// GetTasks は GET /projects/{projectID}/tasks に対応するハンドラ
func (h *handler) GetTasks(_ context.Context, _ ogen.GetTasksParams) (ogen.GetTasksRes, error) {
	return &ogen.GetTasksNotImplemented{}, nil
	//u, err := h.repository.GetUserByToken(ctx, token)
	//if err != nil {
	//	logging.Errorf("getUserByToken failed: %v", err)
	//	return &ogen.GetTasksUnauthorized{}, nil
	//}
	//
	//ts, err := h.repository.GetTasksByUserID(ctx, u.ID)
	//if err != nil {
	//	logging.Errorf("getTasksByUserID failed: %v", err)
	//	return &ogen.GetTasksBadRequest{}, nil
	//}
	//
	//return &ogen.Tasks{Tasks: newTasksResponse(ts)}, nil
}

// PatchTask は PATCH /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *handler) PatchTask(_ context.Context, _ *ogen.PatchTaskReq, _ ogen.PatchTaskParams) (ogen.PatchTaskRes, error) {
	return &ogen.PatchTaskNotImplemented{}, nil
	//u, err := h.repository.GetUserByToken(ctx, token)
	//if err != nil {
	//	logging.Errorf("getUserByToken failed: %v", err)
	//	return &ogen.PatchTaskUnauthorized{}, nil
	//}
	//
	//t, err := h.repository.GetTaskByID(ctx, params.TaskID)
	//if err != nil {
	//	logging.Errorf("getTaskByID failed: %v", err)
	//	return &ogen.PatchTaskBadRequest{}, nil
	//}
	//
	//if t.UserID != u.ID {
	//	logging.Errorf("t.UserID != User.ID")
	//	return &ogen.PatchTaskNotFound{}, nil
	//}
	//
	//if req.IsCompleted.IsSet() {
	//	if req.IsCompleted.Value {
	//		now := time.Now()
	//		if err := h.repository.UpdateTask(ctx, params.TaskID, &now); err != nil {
	//			logging.Errorf("updateTask failed: %v", err)
	//			// TODO: InternalServerError の方が望ましい
	//			return &ogen.PatchTaskBadRequest{}, nil
	//		}
	//	} else {
	//		if err := h.repository.UpdateTask(ctx, params.TaskID, nil); err != nil {
	//			logging.Errorf("updateTask failed: %v", err)
	//			// TODO: InternalServerError の方が望ましい
	//			return &ogen.PatchTaskBadRequest{}, nil
	//		}
	//	}
	//}
	//
	//resp := newTaskResponse(t)
	//return &resp, nil
}

// DeleteTask は DELETE /projects/{projectID}/tasks/{taskID} に対応するハンドラ
func (h *handler) DeleteTask(_ context.Context, _ ogen.DeleteTaskParams) (ogen.DeleteTaskRes, error) {
	return &ogen.DeleteTaskNotImplemented{}, nil
	//u, err := h.repository.GetUserByToken(ctx, token)
	//if err != nil {
	//	logging.Errorf("getUserByToken failed: %v", err)
	//	return &ogen.DeleteTaskUnauthorized{}, nil
	//}
	//
	//t, err := h.repository.GetTaskByID(ctx, params.TaskID)
	//if err != nil {
	//	logging.Errorf("getTaskByID failed: %v", err)
	//	return &ogen.DeleteTaskBadRequest{}, nil
	//}
	//
	//if t.UserID != u.ID {
	//	logging.Errorf("t.UserID != u.ID")
	//	return &ogen.DeleteTaskNotFound{}, nil
	//}
	//
	//if err := h.repository.DeleteTask(ctx, t.ID); err != nil {
	//	logging.Errorf("destroyTask failed: %v", err)
	//	return &ogen.DeleteTaskBadRequest{}, nil
	//}
	//
	//return &ogen.DeleteTaskNoContent{}, nil
}
