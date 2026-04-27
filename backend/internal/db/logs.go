package db

import (
	"database/sql"
	"time"

	"mini-brimble/backend/internal/models"
)

type LogModel struct {
	DB *sql.DB
}

func (m *LogModel) Insert(log *models.Log) error {
	var maxSeq int
	err := m.DB.QueryRow(
		`SELECT COALESCE(MAX(sequence), 0) FROM logs WHERE deployment_id = ?`,
		log.DeploymentId,
	).Scan(&maxSeq)
	if err != nil {
		return err
	}

	log.Sequence = maxSeq + 1
	log.CreatedAt = time.Now()

	_, err = m.DB.Exec(
		`INSERT INTO logs (deployment_id, line, sequence, created_at) VALUES (?, ?, ?, ?)`,
		log.DeploymentId, log.Line, log.Sequence, log.CreatedAt,
	)
	return err
}

func (m *LogModel) GetByDeploymentID(deploymentID string) ([]*models.Log, error) {
	rows, err := m.DB.Query(
		`SELECT id, deployment_id, line, sequence, created_at FROM logs WHERE deployment_id = ? ORDER BY sequence ASC`,
		deploymentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.Log
	for rows.Next() {
		l := &models.Log{}
		err := rows.Scan(&l.Id, &l.DeploymentId, &l.Line, &l.Sequence, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}
