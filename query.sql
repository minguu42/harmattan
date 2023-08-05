-- name: GetUserByAPIKey :one
SELECT * FROM users
WHERE api_key = ?;

-- name: CreateProject :exec
INSERT INTO projects (id, user_id, name, color, is_archived, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetProjectsByUserID :many
SELECT * FROM projects
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetProjectByID :one
SELECT * FROM projects
WHERE id = ?;

-- name: UpdateProject :exec
UPDATE projects SET name = ?, color = ?
WHERE id = ?;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = ?;

-- name: CreateTask :exec
INSERT INTO tasks (id, project_id, title, content, priority, due_on, completed_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTasksByProjectID :many
SELECT * FROM tasks
WHERE project_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: GetTaskByID :one
SELECT * FROM tasks
WHERE id = ?;

-- name: UpdateTask :exec
UPDATE tasks SET title = ?, content = ?, priority = ?, due_on = ?
WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = ?;
