package types

type Services struct {
	Instances    []*Service `yaml:"instances"`
	LastInstance *Service   `yaml:"-"`
}

type Service struct {
	Name      *string `yaml:"name"`
	FullName  *string `yaml:"full_name"`
	ProjectID *int64  `yaml:"gitlab_id"`
	Tag       *string `yaml:"tag"`
	Ports     []*Port `yaml:"ports"`
	LastPort  *Port   `yaml:"-"`
}

type Port struct {
	Number *int64 `yaml:"number"`
}

type OtherServices struct {
	Instances    []*Service `yaml:"instances"`
	LastInstance *Service   `yaml:"-"`
}
