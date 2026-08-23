-- +goose Up
CREATE TABLE project_members (
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (project_id, user_id),
    CONSTRAINT project_members_role_check CHECK (role IN ('OWNER', 'ADMIN', 'MEMBER'))
);

-- Critical for resource authorization: "is this user a member of this project?"
CREATE INDEX project_members_user_project_idx ON project_members(user_id, project_id);

-- Supports member listing for a project.
CREATE INDEX project_members_project_role_idx ON project_members(project_id, role);

-- +goose Down
DROP TABLE project_members;
