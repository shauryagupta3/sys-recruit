package database

import "context"

func (s *service) ApplyToJob(JobId int, UserID int) error {
	_, err := s.db.Query(context.Background(), "INSERT INTO job_profiles (job_id,profile_id) VALUES ($1,$2)", JobId, UserID)
	if err != nil {
		return err
	}
	return nil
}
