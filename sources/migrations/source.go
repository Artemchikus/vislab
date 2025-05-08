package migrations

import (
	"context"
	"vislab/config"
	"vislab/sources/migrations/types"
)

type Source struct {
	parser *Parser
	weight int64
}

func NewSource(config *config.MigrationSourceConfig) (*Source, error) {
	parser, err := NewParser()
	if err != nil {
		return nil, err
	}

	s := &Source{
		parser: parser,
		weight: config.Weight,
	}

	return s, nil
}

func (s *Source) GetData(ctx context.Context, in []byte, out *types.All) error {
	if err := s.parser.Parse(in, out); err != nil {
		return err
	}

	return nil
}

func (s *Source) Weight() int64 {
	return s.weight
}
