package entity

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
			args: args{user: &User{ID: 1}, project: &Project{UserID: 1}},
			want: true,
		},
		{
			name: "ユーザはプロジェクトを所有していない",
			args: args{user: &User{ID: 1}, project: &Project{UserID: 2}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.user.HasProject(tt.args.project); got != tt.want {
				t.Errorf("got: %t, want: %t", got, tt.want)
			}
		})
	}
}
