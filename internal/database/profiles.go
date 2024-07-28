package database

import (
	"context"
	"recruit-sys/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *service) CreateProfile(u *models.Profile) error {
	_, err := s.db.Query(context.Background(), "INSERT INTO profiles (user_id,resume_file_address,skills,education,experience,name,email,phone) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)", u.UserID, u.ResumeFileAddress, u.Skills, u.Education, u.Experience, u.Name, u.Email, u.Phone)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SelectProfileById(id float64) (models.Profile, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM profiles where user_id=$1", id)
	if err != nil {
		return models.Profile{}, err
	}

	profile, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.Profile])
	if err != nil {
		return models.Profile{}, err
	}

	// rows, err = s.db.Query(context.Background(), "SELECT * FROM users where id=$1", id)
	// if err != nil {
	// 	return models.Profile{}, err
	// }

	// user, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[models.User])
	// if err != nil {
	// 	return models.Profile{}, err
	// }

	// profile.User = user

	return profile, err
}

func (s *service) SelectAllProfiles() ([]models.Profile, error) {
	rows, err := s.db.Query(context.Background(), "SELECT * FROM profiles")
	if err != nil {
		return []models.Profile{}, err
	}

	profiles, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[models.Profile])
	if err != nil {
		return []models.Profile{}, err
	}


	return profiles, err
}
