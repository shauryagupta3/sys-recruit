package database

import (
	"context"
	"recruit-sys/internal/models"
)

func (s *service) CreateJob(u *models.Job) error {
	err := s.db.QueryRow(context.Background(), "INSERT INTO jobs (title,description,company_name,posted_by_id) VALUES ($1,$2,$3,$4) RETURNING id", u.Title, u.Description, u.CompanyName, u.PostedByID).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}
