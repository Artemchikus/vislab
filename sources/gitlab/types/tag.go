package types

import "time"

type (
	Tag struct {
		Name      *string    `json:"name"`
		Target    *string    `json:"target"`
		CreatedAt *time.Time `json:"created_at,omitempty"`
	}
	CompareResult struct {
		Commits        []*Commit `json:"commits"`
		Diffs          []*Diff   `json:"diffs"`
		CompareTimeout *bool     `json:"compare_timeout"`
		CompareSameRef *bool     `json:"compare_same_ref"`
	}

	Diff struct {
		OldPath     *string `json:"old_path"`
		NewPath     *string `json:"new_path"`
		NewFile     *bool   `json:"new_file"`
		RenamedFile *bool   `json:"renamed_file"`
		DeletedFile *bool   `json:"deleted_file"`
		Diff        *string `json:"diff"`
	}
)
