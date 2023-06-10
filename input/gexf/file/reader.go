package file

import (
	"graph-analyzer/data-collector/input/gexf/internal"
	"graph-analyzer/data-collector/repository"
	"os"

	log "github.com/sirupsen/logrus"
)

func Reader(repo repository.GraphRepository, fileName string, networkName string) {
	fileContent := readGexfFile(fileName)
	gexfContent, err := internal.UnmarshalGexf(*fileContent)
	if err != nil {
		log.Fatalf("Failed to unmarshal GEXF: %s", err)
	}

	internal.GraphWorker(repo, gexfContent, networkName)
}

func readGexfFile(filename string) *[]byte {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading file from file system: %s", err)
	}

	return &data
}
