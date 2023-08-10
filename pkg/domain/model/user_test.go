package model

import "testing"

func TestUser_HasProject(t *testing.T) {
	type args struct {
		user    *User
		project *Project
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ユーザはプロジェクトを所有している",
			args: args{
				user:    &User{ID: "01DXF6DT000000000000000000"},
				project: &Project{UserID: "01DXF6DT000000000000000000"},
			},
			want: true,
		},
		{
			name: "ユーザはプロジェクトを所有していない",
			args: args{
				user:    &User{ID: "01DXF6DT000000000000000000"},
				project: &Project{UserID: "01DXF6DT000000000000000001"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.user.HasProject(tt.args.project); tt.want != got {
				t.Errorf("user.HasProject want %t, but got %t", tt.want, got)
			}
		})
	}
}
