package types

import "vislab/libs/check"

const (
	ServiceClass     NodeClass = "Service"
	ServicePortClass NodeClass = "ServicePort"
)

type Service struct {
	UID         *string
	Name        *string
	FullName    *string `yaml:"full_name"`
	Link        *string
	Group       *string
	MainBranch  *string
	LatestTag   *string
	Language    *string
	Description *string
	Status      *string
}

func (s Service) Equal(other *Service) bool {
	return check.ComparePointers(s.Name, other.Name) &&
		check.ComparePointers(s.Link, other.Link) &&
		check.ComparePointers(s.Group, other.Group) &&
		check.ComparePointers(s.FullName, other.FullName) &&
		check.ComparePointers(s.MainBranch, other.MainBranch) &&
		check.ComparePointers(s.LatestTag, other.LatestTag) &&
		check.ComparePointers(s.Language, other.Language) &&
		check.ComparePointers(s.Description, other.Description) &&
		check.ComparePointers(s.Status, other.Status)
}

type ServicePort struct {
	UID    *string
	Number *int64
}

func (s ServicePort) Equal(other *ServicePort) bool {
	return check.ComparePointers(s.Number, other.Number)
}
