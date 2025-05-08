package types

import "time"

type (
	ListOptions struct {
		Page    int64 `url:"page,omitempty" json:"page,omitempty"`
		PerPage int64 `url:"per_page,omitempty" json:"per_page,omitempty"`
	}

	ListBranchesOptions struct {
		ListOptions
		Search *string `url:"search,omitempty" json:"search,omitempty"`
		Regex  *string `url:"regex,omitempty" json:"regex,omitempty"`
	}

	ListCommitsOptions struct {
		ListOptions
		RefName *string    `url:"ref_name,omitempty" json:"ref_name,omitempty"`
		Since   *time.Time `url:"since,omitempty" json:"since,omitempty"`
		Until   *time.Time `url:"until,omitempty" json:"until,omitempty"`
		Path    *string    `url:"path,omitempty" json:"path,omitempty"`
	}

	ListContributorsOptions struct {
		ListOptions
		OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
		Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
	}

	ListFilesOptions struct {
		ListOptions
		Ref       *string `url:"ref,omitempty" json:"ref,omitempty"`
		Path      *string `url:"path,omitempty" json:"path,omitempty"`
		Recursive *bool   `url:"recursive,omitempty" json:"recursive,omitempty"`
	}

	ListGroupsOptions struct {
		ListOptions
		Search  *string `url:"search,omitempty" json:"search,omitempty"`
		OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
		Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
	}

	ListMergeRequestsOptions struct {
		ListOptions
		State         *string    `url:"state,omitempty" json:"state,omitempty"`
		SourceBranch  *string    `url:"source_branch,omitempty" json:"source_branch,omitempty"`
		Approved      *bool      `url:"approved,omitempty" json:"approved,omitempty"`
		CreatedAfter  *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
		CreatedBefore *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
		OrderBy       *string    `url:"order_by,omitempty" json:"order_by,omitempty"`
		Search        *string    `url:"search,omitempty" json:"search,omitempty"`
		Sort          *string    `url:"sort,omitempty" json:"sort,omitempty"`
		TargetBranch  *string    `url:"target_branch,omitempty" json:"target_branch,omitempty"`
		UpdatedAfter  *time.Time `url:"updated_after,omitempty" json:"updated_after,omitempty"`
		UpdatedBefore *time.Time `url:"updated_before,omitempty" json:"updated_before,omitempty"`
		Wip           *string    `url:"wip,omitempty" json:"wip,omitempty"`
	}

	ListProjectsOptions struct {
		ListOptions
		Archived           *bool      `url:"archived,omitempty" json:"archived,omitempty"`
		LastActivityAfter  *time.Time `url:"last_activity_after,omitempty" json:"last_activity_after,omitempty"`
		LastActivityBefore *time.Time `url:"last_activity_before,omitempty" json:"last_activity_before,omitempty"`
		Search             *string    `url:"search,omitempty" json:"search,omitempty"`
		OrderBy            *string    `url:"order_by,omitempty" json:"order_by,omitempty"`
		Sort               *string    `url:"sort,omitempty" json:"sort,omitempty"`
	}

	SearchOptions struct {
		ListOptions
		Search  *string `url:"search,omitempty" json:"search,omitempty"`
		Scope   *string `url:"scope,omitempty" json:"scope,omitempty"`
		OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
		Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
		State   *string `url:"state,omitempty" json:"state,omitempty"`
		Ref     *string `url:"ref,omitempty" json:"ref,omitempty"`
	}

	ListTagsOptions struct {
		ListOptions
		Search  *string `url:"search,omitempty" json:"search,omitempty"`
		OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
		Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
	}

	CompareOptions struct {
		ListOptions
		From *string `url:"from,omitempty" json:"from,omitempty"`
		To   *string `url:"to,omitempty" json:"to,omitempty"`
	}
)
