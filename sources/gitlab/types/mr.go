package types

import "time"

type (
	MergeRequest struct {
		ID                       int64      `json:"id"`
		ProjectId                int64      `json:"project_id"`
		Title                    string     `json:"title"`
		ClosedAt                 *time.Time `json:"closed_at"`
		CreatedAt                *time.Time `json:"created_at"`
		MergedAt                 *time.Time `json:"merged_at"`
		Description              string     `json:"description"`
		Draft                    bool       `json:"draft"`
		ShouldRemoveSourceBranch bool       `json:"should_remove_source_branch"`
		CommitSHA                string     `json:"merge_commit_sha"`
		SHA                      string     `json:"sha"`
		SourceBranch             string     `json:"source_branch"`
		State                    string     `json:"state"`
		TargetBranch             string     `json:"target_branch"`
		UpdatedAt                *time.Time `json:"updated_at"`
		Author                   *MRAuthor  `json:"author"`
	}
	MRAuthor struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
		State    string `json:"state"`
		WebURL   string `json:"web_url"`
	}
)
