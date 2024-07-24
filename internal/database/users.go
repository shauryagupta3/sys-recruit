package database

import (
	"context"
	"fmt"
	"recruit-sys/internal/models"
	"strings"
)

func (s *service) CreateUser(u *models.User) error {
	err := s.db.QueryRow(context.Background(), "INSERT INTO users (name,email,address,user_type,password_hash,profile_headline) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id", u.Name, u.Email, u.Address, u.UserType, u.PasswordHash, u.ProfileHeadline).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SelectUserWhereMail(email string) (models.User, error) {
	var user models.User

	query := `
	SELECT id, name, email, address, user_type, password_hash, profile_headline, created_at
	FROM users
	WHERE email = $1`

	err := s.db.QueryRow(context.Background(), query, strings.ToLower(strings.TrimSpace(email))).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.UserType,
		&user.PasswordHash,
		&user.ProfileHeadline,
		&user.CreatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *service) SelectUserWhereID(id float64) (models.User, error) {
	var user models.User

	query := `
	SELECT id, name, email, address, user_type, password_hash, profile_headline, created_at
	FROM users
	WHERE id = $1`

	err := s.db.QueryRow(context.Background(), query,id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Address,
		&user.UserType,
		&user.PasswordHash,
		&user.ProfileHeadline,
		&user.CreatedAt,
	)

	if err != nil {
		fmt.Println(err)
		return models.User{}, err
	}

	return user, nil
}
