package model

import "testing"

func TestProject_ContainsTask(t *testing.T) {
	type args struct {
		project *Project
		task    *Task
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "タスクがプロジェクトに含まれている",
			args: args{
				project: &Project{ID: "01DXF6DT000000000000000000"},
				task:    &Task{ProjectID: "01DXF6DT000000000000000000"},
			},
			want: true,
		},
		{
			name: "タスクがプロジェクトに含まれていない",
			args: args{
				project: &Project{ID: "01DXF6DT000000000000000000"},
				task:    &Task{ProjectID: "01DXF6DT000000000000000001"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.project.ContainsTask(tt.args.task); tt.want != got {
				t.Errorf("project.ContainsTask want %t, but got %t", tt.want, got)
			}
		})
	}
}
