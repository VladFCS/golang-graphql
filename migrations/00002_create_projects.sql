-- +goose Up
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT projects_name_not_blank CHECK (length(trim(name)) > 0)
);

-- Supports listing projects created by a user and auditing ownership.
CREATE INDEX projects_created_by_idx ON projects(created_by);

-- +goose Down
DROP TABLE projects;
