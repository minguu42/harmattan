package entity

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
			name: "タスクはプロジェクトに含まれている",
			args: args{project: &Project{ID: 1}, task: &Task{ProjectID: 1}},
			want: true,
		},
		{
			name: "タスクはプロジェクトに含まれていない",
			args: args{project: &Project{ID: 1}, task: &Task{ProjectID: 2}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.project.ContainsTask(tt.args.task); tt.want != got {
				t.Errorf("project.ContainsTask() want %t, but %t", tt.want, got)
			}
		})
	}
}
