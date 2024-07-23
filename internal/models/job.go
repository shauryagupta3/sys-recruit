package models

import "time"

type Job struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	PostedOn          time.Time `json:"posted_on"`
	TotalApplications int       `json:"total_applications"`
	CompanyName       string    `json:"company_name"`
	PostedByID        uint      `json:"posted_by_id"`
	PostedBy          User      `json:"user"`
	Profiles          []Profile `json:"profiles"`
}
