package input

import (
	"graph-analyzer/data-collector/input/gexf/listener"
	"graph-analyzer/data-collector/repository"

	log "github.com/sirupsen/logrus"
)

type GexfListener struct {
	Port int64
	Host string
}

func (s *GexfListener) read(repo repository.GraphRepository, networkName string) {
	log.Infoln("Using GEXF gRPC listener")
	listener.GrpcServer(repo, s.Host, s.Port)
}
