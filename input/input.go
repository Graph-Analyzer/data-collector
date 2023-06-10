package input

import (
	"graph-analyzer/data-collector/repository"
)

type Input struct {
	Strategy
	NetworkName string
}

func InitInput(s Strategy, networkName string) *Input {
	return &Input{
		s,
		networkName,
	}
}

func (i *Input) Read(repo repository.GraphRepository) {
	i.read(repo, i.NetworkName)
}
