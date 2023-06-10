package jalapeno

import (
	"context"
	"io"
	"sync"

	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	log "github.com/sirupsen/logrus"
)

type EdgeEvent struct {
	Event *jagw.LsNodeEdgeEvent
	Link  *jagw.LsLink
}

func SubscribeService(
	wg *sync.WaitGroup,
	nodeChannel chan *jagw.LsNodeEvent,
	edgeChannel chan EdgeEvent,
	subscribeUrl string,
	requestUrl string,
) {
	defer wg.Done()

	subscribeConnection := newConnection(subscribeUrl)

	requestConnection := newConnection(requestUrl)
	subscribeClient := jagw.NewSubscriptionServiceClient(subscribeConnection)
	requestClient := jagw.NewRequestServiceClient(requestConnection)

	go subscribeLsNodes(subscribeClient, nodeChannel)
	go subscribeLsNodeEdges(subscribeClient, requestClient, edgeChannel)
}

func subscribeLsNodes(client jagw.SubscriptionServiceClient, ch chan *jagw.LsNodeEvent) {
	subscription := &jagw.TopologySubscription{
		Keys:       []string{},
		Properties: []string{"Key", "Id", "Name"},
	}

	log.Infoln("Subscribe to LsNodes")

	stream, err := client.SubscribeToLsNodes(context.Background(), subscription)
	if err != nil {
		log.Fatalf("Error when calling SubscribeToLsNodes on subscribe service: %s", err)
	}

	go func() {
		for {
			log.Debugln("Receiving gRPC event for nodes")
			nodeEvent, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reading from node stream %e", err)
			}

			ch <- nodeEvent
		}
	}()
}

func subscribeLsNodeEdges(subscribeClient jagw.SubscriptionServiceClient, requestClient jagw.RequestServiceClient, ch chan EdgeEvent) {
	subscription := &jagw.TopologySubscription{
		Keys:       []string{},
		Properties: []string{"Key", "Id", "From", "To", "Link"},
	}

	log.Infoln("Subscribe to LsNodesEdges")

	stream, err := subscribeClient.SubscribeToLsNodeEdges(context.Background(), subscription)
	if err != nil {
		log.Fatalf("Error when calling SubscribeToLsNodeEdges on subscribe service: %s", err)
	}

	for {
		log.Debugln("Receiving gRPC event for edges")
		nodeEdgeEvent, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading from edge stream %e", err)
		}

		edgeEvent := EdgeEvent{
			Event: nodeEdgeEvent,
			Link:  nil,
		}

		if *nodeEdgeEvent.Action == "add" {
			link := *nodeEdgeEvent.LsNodeEdge.Link
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

			edgeEvent.Link = linkResponse.LsLinks[0]
		}

		ch <- edgeEvent
	}
}
