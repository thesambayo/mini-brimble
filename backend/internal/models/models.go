package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

// Define Status as a new type based on string
type Status string

// Define the available options as constants of that type
const (
	StatusPending   Status = "pending"
	StatusBuilding  Status = "building"
	StatusDeploying Status = "deploying"
	StatusFailed    Status = "failed"
	StatusRunning   Status = "running"
)

type Deployment struct {
	Id         string    `json:"id"`
	SourceType string    `json:"source_type"`
	Source     string    `json:"source"`
	Status     Status    `json:"status"`
	ImageTag   *string   `json:"image_tag"`
	DeployUrl  *string   `json:"deploy_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Log struct {
	Id           int32     `json:"id"`
	DeploymentId string    `json:"deployment_id"`
	Line         string    `json:"line"`
	Sequence     int       `json:"sequence"`
	CreatedAt    time.Time `json:"created_at"`
}
