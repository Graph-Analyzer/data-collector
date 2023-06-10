package internal

import (
	"fmt"
	"graph-analyzer/data-collector/repository"
	"graph-analyzer/data-collector/repository/models"

	log "github.com/sirupsen/logrus"
	"gonum.org/v1/gonum/graph/formats/gexf12"
)

const (
	undirectedGraphType string = "undirected"
)

func GraphWorker(repo repository.GraphRepository, gexfContent *gexf12.Content, networkName string) {
	graphType := gexfContent.Graph.DefaultEdgeType

	repo.CreateNodeIndex()
	repo.CreateEdgeIndex()

	repo.AddNodes(parseNodes(&gexfContent.Graph.Nodes))
	repo.AddEdges(parseEdges(&gexfContent.Graph.Edges, graphType))

	repo.CreateGraph(networkName)

	log.Infoln("Finished importing GEXF file")
}

func parseNodes(gexfNodes *gexf12.Nodes) *[]models.Node {
	var nodes []models.Node

	for _, importNode := range gexfNodes.Nodes {
		id := importNode.ID
		if id == "" {
			log.Fatal("Node has no id - required by GEXF Schema")
		}

		label := importNode.Label

		nodes = append(nodes, models.Node{
			RouterID:  id,
			RouterKey: id,
			Label:     label,
		})
	}
	log.Infof("Parsed %d nodes from input file", len(nodes))

	return &nodes
}

func parseEdges(gexfEdges *gexf12.Edges, graphType string) *[]models.Edge {
	var edges []models.Edge

	for _, importEdges := range gexfEdges.Edges {
		id := importEdges.ID
		if id == "" {
			log.Fatal("Edge has no ID - required by GEXF Schema")
		}

		source := importEdges.Source
		if source == "" {
			log.Fatal("Edge has no source - required by GEXF Schema")
		}

		target := importEdges.Target
		if target == "" {
			log.Fatal("Edge has no target - required by GEXF Schema")
		}

		weight := importEdges.Weight

		edges = append(edges, models.Edge{
			EdgeID:       id,
			EdgeKey:      id,
			FromRouterID: source,
			ToRouterID:   target,
			Weight:       weight,
			Label:        models.EdgeLabelConnectsTo,
		})

		if graphType == undirectedGraphType {
			id = fmt.Sprintf("%s-%s", id, "rev")

			edges = append(edges, models.Edge{
				EdgeID:       id,
				EdgeKey:      id,
				FromRouterID: target,
				ToRouterID:   source,
				Weight:       weight,
				Label:        models.EdgeLabelConnectsTo,
			})
		}
	}
	log.Infof("Parsed %d edges from input file", len(edges))

	return &edges
}
