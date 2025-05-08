package types

type (
	File struct {
		Name         string `json:"file_name"`
		Path         string `json:"file_path"`
		Ref          string `json:"ref"`
		Content      string `json:"content"`
		BlobId       string `json:"blob_id"`
		CommitId     string `json:"commit_id"`
		LastCommitId string `json:"last_commit_id"`
	}

	ListFile struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		Path string `json:"path"`
		Mode string `json:"mode"`
	}
)
