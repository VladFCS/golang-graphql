package project

import "time"

type Project struct {
	ID          string
	Name        string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateProjectInput struct {
	Name        string
	Description string
	CreatedBy   string
}
