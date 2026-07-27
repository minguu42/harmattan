package domain

import "time"

// MaxProjectsPerUser は1ユーザが作成できるプロジェクト数の上限
const MaxProjectsPerUser = 100

type ProjectID string

type Project struct {
	ID         ProjectID
	UserID     UserID
	Name       string
	Color      ProjectColor
	IsArchived bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Projects []Project

type ProjectColor string

const (
	ProjectColorBlue    ProjectColor = "blue"
	ProjectColorBrown   ProjectColor = "brown"
	ProjectColorDefault ProjectColor = "default"
	ProjectColorGray    ProjectColor = "gray"
	ProjectColorGreen   ProjectColor = "green"
	ProjectColorOrange  ProjectColor = "orange"
	ProjectColorPink    ProjectColor = "pink"
	ProjectColorPurple  ProjectColor = "purple"
	ProjectColorRed     ProjectColor = "red"
	ProjectColorYellow  ProjectColor = "yellow"
)
