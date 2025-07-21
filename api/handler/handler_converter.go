package handler

import (
	"time"

	"github.com/minguu42/harmattan/internal/domain"
	"github.com/minguu42/harmattan/internal/openapi"
	"github.com/minguu42/harmattan/lib/pointers"
)

func convertOptBool(v openapi.OptBool) *bool {
	if v.Set {
		return &v.Value
	}
	return nil
}

func convertOptInt(o openapi.OptInt) *int {
	if o.Set {
		return &o.Value
	}
	return nil
}

func convertOptString(v openapi.OptString) *string {
	if v.Set {
		return &v.Value
	}
	return nil
}

func convertOptColorString(v openapi.OptUpdateProjectReqColor) *domain.ProjectColor {
	if v.Set {
		return pointers.Ref(domain.ProjectColor(v.Value))
	}
	return nil
}

func convertDatePtr(v *time.Time) openapi.OptDate {
	if v != nil {
		return openapi.OptDate{Value: *v, Set: true}
	}
	return openapi.OptDate{}
}

func convertOptDateTime(v openapi.OptDateTime) *time.Time {
	if v.Set {
		return &v.Value
	}
	return nil
}

func convertDateTimePtr(v *time.Time) openapi.OptDateTime {
	if v != nil {
		return openapi.OptDateTime{Value: *v, Set: true}
	}
	return openapi.OptDateTime{}
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
		CompletedAt: convertDateTimePtr(s.CompletedAt),
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

func convertTask(task *domain.Task) *openapi.Task {
	return &openapi.Task{
		ID:          string(task.ID),
		ProjectID:   string(task.ProjectID),
		Name:        task.Name,
		Content:     task.Content,
		Priority:    task.Priority,
		DueOn:       convertDatePtr(task.DueOn),
		CompletedAt: convertDateTimePtr(task.CompletedAt),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Steps:       convertSteps(task.Steps),
		Tags:        convertTags(task.Tags),
	}
}

func convertTasks(tasks domain.Tasks) []openapi.Task {
	ts := make([]openapi.Task, 0, len(tasks))
	for _, t := range tasks {
		ts = append(ts, *convertTask(&t))
	}
	return ts
}
