package internal

import (
	"encoding/xml"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph/formats/gexf12"
)

func UnmarshalGexf(data []byte) (*gexf12.Content, error) {
	var graph gexf12.Content

	err := xml.Unmarshal(data, &graph)
	if err != nil {
		log.Infoln("Error unmarshal GEXF content")
		return nil, err
	}

	return &graph, nil
}
