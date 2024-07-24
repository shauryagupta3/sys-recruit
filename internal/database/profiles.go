package database

import (
	"context"
	"recruit-sys/internal/models"
)

func (s *service) CreateProfile(u *models.Profile) error {
	_, err := s.db.Query(context.Background(), "INSERT INTO profiles (user_id,resume_file_address,skills,education,experience,name,email,phone) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)", u.UserID, u.ResumeFileAddress, u.Skills, u.Education, u.Experience, u.Name, u.Email, u.Phone)
	if err != nil {
		return err
	}
	return nil
}
