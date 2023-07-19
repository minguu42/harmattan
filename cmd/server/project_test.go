package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/minguu42/mtasks/gen/mock"
)

func TestProject(t *testing.T) {
	run(t, []test{
		{
			id:            "createProject",
			method:        http.MethodPost,
			path:          "/projects",
			body:          strings.NewReader(`{"name": "新プロジェクト", "color": "#1A2B3C"}`),
			statusCode:    http.StatusCreated,
			needsRollback: true,
			prepareMockFn: func(t *testing.T) {
				m := mock.NewMockIDGenerator(gomock.NewController(t))
				m.EXPECT().Generate().Return("01DXF6DT000000000000000002")
				tdb.SetIDGenerator(m)
			},
		},
		{
			id:         "listProject",
			method:     http.MethodGet,
			path:       "/projects?sort=createdAt",
			statusCode: http.StatusOK,
		},
		{
			id:            "updateProject",
			method:        http.MethodPatch,
			path:          "/projects/01DXF6DT000000000000000000",
			body:          strings.NewReader(`{"name": "新プロジェクト1"}`),
			statusCode:    http.StatusOK,
			needsRollback: true,
		},
		{
			id:            "deleteProject",
			method:        http.MethodDelete,
			path:          "/projects/01DXF6DT000000000000000000",
			statusCode:    http.StatusNoContent,
			needsRollback: true,
		},
	})
}
