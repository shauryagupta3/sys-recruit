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
	address VARCHAR(100),
	user_type VARCHAR(50) NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);
`
const createProfilesTable = `
CREATE TABLE IF NOT EXISTS profiles (
    user_id INTEGER PRIMARY KEY,
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
    title VARCHAR(255) NOT NULL,
    description TEXT,
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
