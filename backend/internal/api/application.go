package api

import (
	"mini-brimble/backend/internal/db"
	"mini-brimble/backend/internal/pipeline"
	"sync"
)

type Application struct {
	Deployments *db.DeploymentModel
	Logs        *db.LogModel
	Pipeline    *pipeline.Pipeline
	LogStreams  map[string]chan string
	StreamsMu   sync.Mutex
}
