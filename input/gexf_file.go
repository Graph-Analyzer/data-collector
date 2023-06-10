package input

import (
	"graph-analyzer/data-collector/input/gexf/file"
	"graph-analyzer/data-collector/repository"

	log "github.com/sirupsen/logrus"
)

type GexfFile struct {
	Filename string
}

func (s *GexfFile) read(repo repository.GraphRepository, networkName string) {
	log.Infof("Using GEXF file %s", s.Filename)
	file.Reader(repo, s.Filename, networkName)
}
