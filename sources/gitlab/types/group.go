package types

type (
	Group struct {
		ID          int64  `json:"id"`
		Name        string `json:"name"`
		Path        string `json:"path"`
		Description string `json:"description"`
		FullName    string `json:"full_name"`
		FullPath    string `json:"full_path"`
		WebURL      string `json:"web_url"`
	}
)
