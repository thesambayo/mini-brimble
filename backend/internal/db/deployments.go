package db

import (
	"database/sql"
	"errors"
	"time"

	"mini-brimble/backend/internal/models"

	"github.com/google/uuid"
)

type DeploymentModel struct {
	DB *sql.DB
}

func (m *DeploymentModel) Create(d *models.Deployment) (string, error) {
	d.Id = uuid.New().String()
	d.Status = models.StatusPending
	now := time.Now()
	d.CreatedAt = now
	d.UpdatedAt = now

	stmt := `INSERT INTO deployments (id, source_type, source, status, image_tag, deploy_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt, d.Id, d.SourceType, d.Source, d.Status, d.ImageTag, d.DeployUrl, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return "", err
	}
	return d.Id, nil
}

func (m *DeploymentModel) Get(id string) (*models.Deployment, error) {
	stmt := `SELECT id, source_type, source, status, image_tag, deploy_url, created_at, updated_at
		FROM deployments WHERE id = ?`

	d := &models.Deployment{}
	err := m.DB.QueryRow(stmt, id).Scan(
		&d.Id,
		&d.SourceType,
		&d.Source,
		&d.Status,
		&d.ImageTag,
		&d.DeployUrl,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}
	return d, nil
}

func (m *DeploymentModel) GetAll() ([]*models.Deployment, error) {
	stmt := `SELECT id, source_type, source, status, image_tag, deploy_url, created_at, updated_at
		FROM deployments ORDER BY created_at DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deployments := []*models.Deployment{}
	for rows.Next() {
		d := &models.Deployment{}
		err := rows.Scan(
			&d.Id,
			&d.SourceType,
			&d.Source,
			&d.Status,
			&d.ImageTag,
			&d.DeployUrl,
			&d.CreatedAt,
			&d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, d)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return deployments, nil
}

func (m *DeploymentModel) UpdateStatus(id string, status models.Status) error {
	stmt := `UPDATE deployments SET status = ?, updated_at = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, status, time.Now(), id)
	return err
}

func (m *DeploymentModel) UpdateImageTag(id string, imageTag string) error {
	stmt := `UPDATE deployments SET image_tag = ?, updated_at = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, imageTag, time.Now(), id)
	return err
}

func (m *DeploymentModel) UpdateDeployURL(id string, deployURL string) error {
	stmt := `UPDATE deployments SET deploy_url = ?, updated_at = ? WHERE id = ?`
	_, err := m.DB.Exec(stmt, deployURL, time.Now(), id)
	return err
}
