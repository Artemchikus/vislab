package types

type (
	SearchResult struct {
		Basename  string `json:"basename"`
		Data      string `json:"data"`
		Path      string `json:"path"`
		ProjectId int64  `json:"project_id"`
		Ref       string `json:"ref"`
		StartLine int64  `json:"start_line"`
	}
	FileQuery struct {
		Path      string
		Filename  string
		Extension string
	}
)
