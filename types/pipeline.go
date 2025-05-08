package types

type Pipeline struct {
	Name *string
	Link *string
	Jobs []*PipelineJob
}

type PipelineJob struct {
	Name *string
	Link *string
}
