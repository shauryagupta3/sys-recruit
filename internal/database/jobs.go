package database

import (
	"context"
	"recruit-sys/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *service) CreateJob(u *models.Job) error {
	err := s.db.QueryRow(context.Background(), "INSERT INTO jobs (title,description,company_name,posted_by_id) VALUES ($1,$2,$3,$4) RETURNING id", u.Title, u.Description, u.CompanyName, u.PostedByID).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SelectAllJobs() ([]models.Job, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM jobs")
	if err != nil {
		return nil, err
	}

	jobs, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Job])
	if err != nil {
		return nil, err
	}
	return jobs, err
}

func (s *service) SelectJobsByID(id int) (models.Job, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM jobs where id=$1", id)
	if err != nil {
		return models.Job{}, err
	}

	jobs, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Job])
	if err != nil {
		return models.Job{}, err
	}
	return jobs, err
}

func (s *service) SelectJobsPostedBy(postedById float64) ([]models.Job, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM jobs where posted_by_id=$1", postedById)
	if err != nil {
		return nil, err
	}

	jobs, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Job])
	if err != nil {
		return nil, err
	}
	return jobs, err
}

func (s *service) SelectJobByIdAdmin(postedById float64,id int) (models.Job, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM jobs where id=$1 AND posted_by_id=$2", id,postedById)
	if err != nil {
		return models.Job{}, err
	}

	user,err := s.SelectUserWhereID(postedById);
	if err!=nil {
		return models.Job{},err
	}

	job, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Job])
	if err != nil {
		return models.Job{}, err
	}

	job.PostedBy = user

	return job, err
}
