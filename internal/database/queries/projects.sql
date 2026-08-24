-- name: CreateProject :one
INSERT INTO projects (
    name,
    description,
    created_by
) VALUES (
    $1,
    $2,
    $3
)
RETURNING
    id,
    name,
    description,
    created_by,
    created_at,
    updated_at;

-- name: ListProjects :many
SELECT
    id,
    name,
    description,
    created_by,
    created_at,
    updated_at
FROM projects
ORDER BY created_at DESC, id DESC;

-- name: GetProjectByID :one
SELECT
    id,
    name,
    description,
    created_by,
    created_at,
    updated_at
FROM projects
WHERE id = $1;

-- name: ListProjectsByUserID :many
SELECT
    p.id,
    p.name,
    p.description,
    p.created_by,
    p.created_at,
    p.updated_at
FROM projects AS p
JOIN project_members AS pm ON pm.project_id = p.id
WHERE pm.user_id = $1
ORDER BY p.created_at DESC, p.id DESC;
