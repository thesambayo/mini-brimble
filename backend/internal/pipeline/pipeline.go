package pipeline

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"mini-brimble/backend/internal/db"
	"mini-brimble/backend/internal/models"
	"net"
	"net/http"
	"os"
	"os/exec"
)

type Pipeline struct {
	Deployments   *db.DeploymentModel
	Logs          *db.LogModel
	WorkspaceDir  string
	CaddyAdminURL string
	DockerNetwork string
}

func (p *Pipeline) Clone(deploymentId, gitUrl string) (string, error) {
	// create the workspace directory
	workspacePath := fmt.Sprintf("%s/%s", p.WorkspaceDir, deploymentId)
	if err := os.MkdirAll(workspacePath, 0755); err != nil {
		return "", fmt.Errorf("error creating workspace: %w", err)
	}

	// clone the repo into created directory
	cmd := exec.Command("git", "clone", gitUrl, workspacePath)
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("error cloning repository: %w", err)
	}

	fmt.Println("cloning done")

	return workspacePath, nil
}

func (p *Pipeline) Build(deploymentId, workspacePath, appName string, logCh chan string) (string, error) {
	imageTag := fmt.Sprintf("%s:%s", appName, deploymentId[:8])

	cmd := exec.Command("railpack", "build", "--name", imageTag, ".")
	cmd.Dir = workspacePath

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("error piping build output: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("error starting build: %w", err)
	}

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

		if err := p.Logs.Insert(&models.Log{
			DeploymentId: deploymentId,
			Line:         line,
		}); err != nil {
			fmt.Printf("failed to insert log: %v\n", err)
		}

		// send to SSE channel
		logCh <- line
	}

	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("build failed: %w", err)
	}

	return imageTag, nil
}

func (p *Pipeline) Run(deploymentId, imageTag string) (int, error) {
	listener, err := net.Listen("tcp", ":0")

	if err != nil {
		return 0, err
	}
	// ":0" tells the OS to assign any free port
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	cmd := exec.Command("docker", "run", "-d",
		"-p", fmt.Sprintf("%d:%d", port, 3000),
		"--name", deploymentId,
		"--network", p.DockerNetwork,
		imageTag,
	)
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("error starting container: %w", err)
	}
	return port, nil
}

func (p *Pipeline) RegisterRoute(deploymentId string) (string, error) {
	subdomain := deploymentId[:8]
	upstreamDial := fmt.Sprintf("%s:%d", deploymentId, 3000)

	body := map[string]any{
		"match": []map[string]any{
			{"host": []string{fmt.Sprintf("%s.localhost", subdomain)}},
		},
		"handle": []map[string]any{
			{
				"handler":   "reverse_proxy",
				"upstreams": []map[string]any{{"dial": upstreamDial}},
			},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("error marshalling route body: %w", err)
	}

	url := fmt.Sprintf("%s/config/apps/http/servers/main/routes", p.CaddyAdminURL)
	resp, err := http.Post(url, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("error calling caddy admin api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("caddy returned non-200 status: %d", resp.StatusCode)
	}

	return fmt.Sprintf("http://%s.localhost", subdomain), nil
}
