package input

import (
	"graph-analyzer/data-collector/repository"
)

type Strategy interface {
	read(repo repository.GraphRepository, networkName string)
}
