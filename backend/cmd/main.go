package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"mini-brimble/backend/internal/api"
	"mini-brimble/backend/internal/db"
	"mini-brimble/backend/internal/pipeline"
)

func main() {
	// env variables
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "data/app.db"
	}

	workspaceDir := os.Getenv("WORKSPACE_DIR")
	if workspaceDir == "" {
		workspaceDir = "/tmp/brimble-workspaces"
	}

	caddyAdminURL := os.Getenv("CADDY_ADMIN_URL")
	if caddyAdminURL == "" {
		caddyAdminURL = "http://localhost:2019"
	}

	dockerNetwork := os.Getenv("DOCKER_NETWORK")
	if dockerNetwork == "" {
		dockerNetwork = "mini-brimble_brimble-net"
	}
	// env ends

	conn, err := db.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Critical error initializing database: %v", err)
	}
	defer conn.Close()

	if err := db.InitSchema(conn); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	app := &api.Application{
		Deployments: &db.DeploymentModel{DB: conn},
		Logs:        &db.LogModel{DB: conn},
		Pipeline: &pipeline.Pipeline{
			Deployments:   &db.DeploymentModel{DB: conn},
			Logs:          &db.LogModel{DB: conn},
			WorkspaceDir:  workspaceDir,
			CaddyAdminURL: caddyAdminURL,
			DockerNetwork: dockerNetwork,
		},
		LogStreams: make(map[string]chan string),
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	mux := app.Routes()

	log.Printf("Starting server on :%s", port)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   0,
		MaxHeaderBytes: 524288,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
