package types

import "time"

type (
	Commit struct {
		ID            string     `json:"id"`
		CommittedDate *time.Time `json:"committed_date"`
		CreatedAt     *time.Time `json:"created_at"`
	}
)
