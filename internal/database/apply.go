package database

import (
	"context"
	"recruit-sys/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *service) ApplyToJob(JobId int, UserID int) error {
	_, err := s.db.Query(context.Background(), "INSERT INTO job_profiles (job_id,profile_id) VALUES ($1,$2)", JobId, UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SelectJobsAppliedBy(userId float64) ([]models.Job, error) {
	rows, err := s.db.Query(context.Background(), "SELECT j.* FROM job_profiles INNER JOIN jobs j ON j.id=job_profiles.job_id WHERE profile_id=$1", userId)
	if err != nil {
		return nil, err
	}
	jobs, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Job])
	if err != nil {
		return nil, err
	}
	return jobs, err
}

func (s *service) SelectProfilesAppliedBy(jobId int) ([]models.Profile, error) {
	rows, err := s.db.Query(context.Background(), "SELECT p.* FROM job_profiles INNER JOIN profiles p ON p.user_id=job_profiles.profile_id WHERE job_id=$1", jobId)
	if err != nil {
		return nil, err
	}
	profiles, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Profile])
	if err != nil {
		return nil, err
	}
	return profiles, err
}