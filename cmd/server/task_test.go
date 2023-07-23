package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/minguu42/mtasks/gen/mock"
)

func TestTask(t *testing.T) {
	run(t, []test{
		{
			id:         "createTask",
			method:     http.MethodPost,
			path:       "/projects/01DXF6DT000000000000000000/tasks",
			body:       strings.NewReader(`{"title": "新タスク", "content": "Hello, 世界!", "priority": 3, "dueOn": "2020-01-02"}`),
			statusCode: http.StatusCreated,
			prepareMockFn: func(t *testing.T) {
				g := mock.NewMockIDGenerator(gomock.NewController(t))
				g.EXPECT().Generate().Return("01DXF6DT000000000000000002")
				tdb.SetIDGenerator(g)
			},
			needsRollback: true,
		},
		{
			id:         "listTasks",
			method:     http.MethodGet,
			path:       "/projects/01DXF6DT000000000000000000/tasks?sort=createdAt",
			statusCode: http.StatusOK,
		},
		{
			id:            "updateTask",
			method:        http.MethodPatch,
			path:          "/projects/01DXF6DT000000000000000000/tasks/01DXF6DT000000000000000000",
			body:          strings.NewReader(`{"isCompleted": true}`),
			statusCode:    http.StatusOK,
			needsRollback: true,
		},
		{
			id:            "deleteTask",
			method:        http.MethodDelete,
			path:          "/projects/01DXF6DT000000000000000000/tasks/01DXF6DT000000000000000000",
			statusCode:    http.StatusNoContent,
			needsRollback: true,
		},
	})
}
