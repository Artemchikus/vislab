package types

type (
	Branch struct {
		Name    string  `json:"name"`
		Commit  *Commit `json:"commit"`
		Merged  bool    `json:"merged"`
		Default bool    `json:"default"`
	}
)
