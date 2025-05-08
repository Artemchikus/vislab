package types

import "time"

type (
	Project struct {
		ID             *int64     `json:"id"`
		Description    *string    `json:"description"`
		DefaultBranch  *string    `json:"default_branch"`
		WebURL         *string    `json:"web_url"`
		Name           *string    `json:"name"`
		PathWithGroup  *string    `json:"path_with_namespace"`
		CreatedAt      *time.Time `json:"created_at,omitempty"`
		UpdatedAt      *time.Time `json:"updated_at,omitempty"`
		LastActivityAt *time.Time `json:"last_activity_at,omitempty"`
		Group          *Namespace `json:"namespace,omitempty"`
		EmptyRepo      *bool      `json:"empty_repo"`
		Archived       *bool      `json:"archived"`
		Owner          *Owner     `json:"owner,omitempty"`
	}

	Namespace struct {
		ID   *int64  `json:"id"`
		Name *string `json:"name"`
		Path *string `json:"path"`
		Kind *string `json:"kind"`
	}

	Owner struct {
		ID       *int64  `json:"id"`
		Username *string `json:"username"`
		Email    *string `json:"email"`
		Name     *string `json:"name"`
		State    *string `json:"state"`
		Note     *string `json:"note"`
		Locked   *bool   `json:"locked"`
	}
)
