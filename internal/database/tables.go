package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	address VARCHAR(100) not null,
	user_type VARCHAR(50) NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
    profile_headline VARCHAR(255) not null,
	created_at TIMESTAMP DEFAULT NOW()
);
`
const createProfilesTable = `
CREATE TABLE IF NOT EXISTS profiles (
    user_id INT UNIQUE,
    resume_file_address VARCHAR(255),
    skills TEXT,
    education TEXT,
    experience TEXT,
    name VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);`

const createJobsTable = `
CREATE TABLE IF NOT EXISTS jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    posted_on TIMESTAMP DEFAULT NOW(),
    total_applications INTEGER DEFAULT 0,
    company_name VARCHAR(255) NOT NULL,
    posted_by_id INTEGER NOT NULL,
    FOREIGN KEY (posted_by_id) REFERENCES users(id) ON DELETE SET NULL
);`

const createJobProfileTable = `
CREATE TABLE IF NOT EXISTS job_profiles (
    job_id INTEGER NOT NULL,
    profile_id INTEGER NOT NULL,
    PRIMARY KEY (job_id, profile_id),
    FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
    FOREIGN KEY (profile_id) REFERENCES profiles(user_id) ON DELETE CASCADE
);`

func CreateTables(ctx context.Context, db *pgxpool.Pool) error {
	tables := []string{createUsersTable, createProfilesTable, createJobsTable, createJobProfileTable}

	for _, table := range tables {
		_, err := db.Exec(ctx, table)
		if err != nil {
			return fmt.Errorf("failed to execute SQL statement: %w", err)
		}

	}

	return nil
}

func DropTables(ctx context.Context, db *pgxpool.Pool) error {
	tables := []string{
		"users",
		"jobs",
		"profiles",
		"job_profiles",
	}

	for _, query := range tables {
		query := fmt.Sprintf("drop table if exists %s cascade", query)
		_, err := db.Exec(ctx, query)
		if err != nil {
			return err
		}
	}
	return nil
}
