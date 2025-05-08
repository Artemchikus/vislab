package types

type Service struct {
	Name        *string
	Link        *string
	Group       *string
	FullName    *string
	MainBranch  *string
	LatestTag   *string
	Language    *string
	Description *string
	Status      *string
	Ports       []*Port
}

type Port struct {
	Number *int64
}
