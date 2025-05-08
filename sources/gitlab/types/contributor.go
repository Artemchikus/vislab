package types

type (
	Contributor struct {
		Name       string `json:"name"`
		Email      string `json:"email"`
		CommitsNum int64  `json:"commits"`
	}
)
