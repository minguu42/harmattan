package handler

import (
	"time"

	"github.com/minguu42/harmattan/internal/api/openapi"
	"github.com/minguu42/harmattan/internal/domain"
)

func ternary[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func convertOptDate(t *time.Time) openapi.OptDate {
	if t != nil {
		return openapi.OptDate{Value: *t, Set: true}
	}
	return openapi.OptDate{}
}

func convertOptDateTime(t *time.Time) openapi.OptDateTime {
	if t != nil {
		return openapi.OptDateTime{Value: *t, Set: true}
	}
	return openapi.OptDateTime{}
}

func convertSlice[T ~string](s []string) []T {
	r := make([]T, 0, len(s))
	for _, e := range s {
		r = append(r, T(e))
	}
	return r
}

func convertProject(project *domain.Project) *openapi.Project {
	return &openapi.Project{
		ID:         string(project.ID),
		Name:       project.Name,
		Color:      openapi.ProjectColor(project.Color),
		IsArchived: project.IsArchived,
		CreatedAt:  project.CreatedAt,
		UpdatedAt:  project.UpdatedAt,
	}
}

func convertProjects(projects domain.Projects) []openapi.Project {
	ps := make([]openapi.Project, 0, len(projects))
	for _, p := range projects {
		ps = append(ps, *convertProject(&p))
	}
	return ps
}

func convertStep(s *domain.Step) *openapi.Step {
	return &openapi.Step{
		ID:          string(s.ID),
		TaskID:      string(s.TaskID),
		Name:        s.Name,
		CompletedAt: convertOptDateTime(s.CompletedAt),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

func convertSteps(steps domain.Steps) []openapi.Step {
	s := make([]openapi.Step, 0, len(steps))
	for _, step := range steps {
		s = append(s, *convertStep(&step))
	}
	return s
}

func convertTag(t *domain.Tag) *openapi.Tag {
	return &openapi.Tag{
		ID:        string(t.ID),
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func convertTags(tags domain.Tags) []openapi.Tag {
	ts := make([]openapi.Tag, 0, len(tags))
	for _, t := range tags {
		ts = append(ts, *convertTag(&t))
	}
	return ts
}

func convertTask(task *domain.Task, tags domain.Tags) *openapi.Task {
	return &openapi.Task{
		ID:          string(task.ID),
		ProjectID:   string(task.ProjectID),
		Name:        task.Name,
		Content:     task.Content,
		Priority:    task.Priority,
		DueOn:       convertOptDate(task.DueOn),
		CompletedAt: convertOptDateTime(task.CompletedAt),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Steps:       convertSteps(task.Steps),
		Tags:        convertTags(tags),
	}
}

func convertTasks(tasks domain.Tasks, tags domain.Tags) []openapi.Task {
	tagByID := tags.TagByID()

	ts := make([]openapi.Task, 0, len(tasks))
	for _, t := range tasks {
		taskTags := make(domain.Tags, 0, len(t.TagIDs))
		for _, id := range t.TagIDs {
			taskTags = append(taskTags, tagByID[id])
		}
		ts = append(ts, *convertTask(&t, taskTags))
	}
	return ts
}
