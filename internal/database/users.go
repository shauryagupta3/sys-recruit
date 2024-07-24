package database

import (
	"context"
	"recruit-sys/internal/models"

	"github.com/jackc/pgx/v5"
)

func (s *service) CreateUser(u *models.User) error {
	err := s.db.QueryRow(context.Background(), "INSERT INTO users (name,email,address,user_type,password_hash,profile_headline) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id", u.Name, u.Email, u.Address, u.UserType, u.PasswordHash, u.ProfileHeadline).Scan(&u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) SelectUserWhereMail(email string) (models.User, error) {
	// var user models.User
	
	query := `
	SELECT id, name, email, address, user_type, password_hash, profile_headline, created_at
	FROM users
	WHERE email = $1`

	// err := s.db.QueryRow(context.Background(), query, strings.ToLower(strings.TrimSpace(email))).Scan(
	// 	&user.ID,
	// 	&user.Name,
	// 	&user.Email,
	// 	&user.Address,
	// 	&user.UserType,
	// 	&user.PasswordHash,
	// 	&user.ProfileHeadline,
	// 	&user.CreatedAt,
	// )

	// if err != nil {
	// 	return models.User{}, err
	// }

	rows,_ := s.db.Query(context.Background(),query,email)
	p,err := pgx.CollectOneRow(rows,pgx.RowToStructByName[models.User])
	if err!=nil {
		return models.User{},err
	}
	return p, nil
}

func (s *service) SelectUserWhereID(id float64) (models.User, error) {
	query := `
	SELECT id, name, email, address, user_type, password_hash, profile_headline, created_at
	FROM users
	WHERE id = $1`

	rows,_ := s.db.Query(context.Background(),query,id)
	p,err := pgx.CollectOneRow(rows,pgx.RowToStructByName[models.User])
	if err!=nil {
		return models.User{},err
	}

	return p, nil
}
