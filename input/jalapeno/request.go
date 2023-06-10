package jalapeno

import (
	"context"
	"graph-analyzer/data-collector/repository"
	"graph-analyzer/data-collector/repository/models"
	"sync"

	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func RequestService(repo repository.GraphRepository, requestUrl string, networkName string) {
	requestConnection := newConnection(requestUrl)
	defer func(requestConnection *grpc.ClientConn) {
		err := requestConnection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(requestConnection)

	requestClient := jagw.NewRequestServiceClient(requestConnection)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	repo.CreateNodeIndex()
	repo.CreateEdgeIndex()

	getLsNodes(wg, requestClient, repo)
	getLsNodeEdges(wg, requestClient, repo)

	wg.Wait()

	repo.CreateGraph(networkName)
}

func getLsNodes(wg *sync.WaitGroup, requestClient jagw.RequestServiceClient, repo repository.GraphRepository) {
	defer wg.Done()
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "Id", "Name"},
	}

	response, err := requestClient.GetLsNodes(context.Background(), request)
	if err != nil {
		log.Fatalf("Error when calling GetLsNodes on request service: %s", err)
	}

	var nodes []models.Node
	for i := 0; i < len(response.LsNodes); i++ {
		nodes = append(nodes, models.Node{
			RouterID:  *response.LsNodes[i].Id,
			RouterKey: *response.LsNodes[i].Key,
			Label:     *response.LsNodes[i].Name,
		})
	}

	log.Infof("Received %d nodes from JAGW getLsNodes", len(nodes))
	repo.AddNodes(&nodes)
}

func getLsNodeEdges(wg *sync.WaitGroup, requestClient jagw.RequestServiceClient, repo repository.GraphRepository) {
	defer wg.Done()
	request := &jagw.TopologyRequest{
		Keys:       []string{},
		Properties: []string{"Key", "Id", "From", "To", "Link"},
	}

	response, err := requestClient.GetLsNodeEdges(context.Background(), request)
	if err != nil {
		log.Fatalf("Error when calling GetLsNodeEdges on request service: %s", err)
	}

	var edges []models.Edge

	for i := 0; i < len(response.LsNodeEdges); i++ {
		link := *response.LsNodeEdges[i].Link
		linkRequest := &jagw.TopologyRequest{
			Keys:       []string{link},
			Properties: []string{"Key", "Id", "IgpMetric"},
		}

		linkResponse, err := requestClient.GetLsLinks(context.Background(), linkRequest)
		if err != nil {
			log.Fatalf("Error when calling GetLsLinks on request service: %s", err)
		}

		if len(linkResponse.LsLinks) != 1 {
			log.Fatal("Unexpected response from LsLinks encountered in request service")
		}

		edges = append(edges, models.Edge{
			EdgeID:       *response.LsNodeEdges[i].Id,
			EdgeKey:      *response.LsNodeEdges[i].Key,
			FromRouterID: *response.LsNodeEdges[i].From,
			ToRouterID:   *response.LsNodeEdges[i].To,
			Weight:       float64(*linkResponse.LsLinks[0].IgpMetric),
			Label:        models.EdgeLabelConnectsTo,
		})
	}

	log.Infof("Received %d edges from JAGW getLsNodeEdges", len(edges))
	repo.AddEdges(&edges)
}
