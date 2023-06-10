package jalapeno

import (
	"graph-analyzer/data-collector/repository"
	"graph-analyzer/data-collector/repository/models"
	"sync"
	"time"

	"github.com/jalapeno-api-gateway/jagw-go/jagw"
	log "github.com/sirupsen/logrus"
)

var waitMutex = &sync.Mutex{}

func ProcessData(repo repository.GraphRepository, nodeChannel chan *jagw.LsNodeEvent, edgeChannel chan EdgeEvent, networkName string) {
	log.Infoln("Start processing node updates")
	go func() {
		for nodeEvent := range nodeChannel {
			switch *nodeEvent.Action {
			case "del":
				log.Infof("Deleting node by key %s", *nodeEvent.Key)
				repo.DeleteNodeByKey(*nodeEvent.Key)
				recreateGraph(repo, networkName)
			case "update":
				log.Infof("Updating node label %s of node %s", *nodeEvent.LsNode.Name, *nodeEvent.LsNode.Id)
				repo.UpdateNodeLabelById(*nodeEvent.LsNode.Id, *nodeEvent.LsNode.Name)
				recreateGraph(repo, networkName)
			case "add":
				node := models.Node{
					RouterID:  *nodeEvent.LsNode.Id,
					RouterKey: *nodeEvent.Key,
					Label:     *nodeEvent.LsNode.Name,
				}

				nodeExists, err := repo.CheckNodeExistsById(*nodeEvent.LsNode.Id)
				if err != nil {
					log.Fatalf("Error checking if node exists: %s", err)
				}

				if !nodeExists {
					log.Infof("Adding node %s", *nodeEvent.LsNode.Name)
					repo.AddNode(node)
					recreateGraph(repo, networkName)
				}
			}
		}
	}()

	log.Infoln("Start processing edge updates")
	go func() {
		for edgeEvent := range edgeChannel {
			switch *edgeEvent.Event.Action {
			case "del":
				log.Infof("Deleting edge by key %s", *edgeEvent.Event.Key)
				repo.DeleteEdgeByKey(*edgeEvent.Event.Key)
				recreateGraph(repo, networkName)
			case "add":
				edge := models.Edge{
					EdgeID:       *edgeEvent.Event.LsNodeEdge.Id,
					EdgeKey:      *edgeEvent.Event.LsNodeEdge.Key,
					FromRouterID: *edgeEvent.Event.LsNodeEdge.From,
					ToRouterID:   *edgeEvent.Event.LsNodeEdge.To,
					Weight:       float64(*edgeEvent.Link.IgpMetric),
					Label:        models.EdgeLabelConnectsTo,
				}

				edgeExists, errExists := repo.CheckEdgeExistsById(*edgeEvent.Event.LsNodeEdge.Id)
				if errExists != nil {
					log.Fatalf("Error checking if edge exists: %s", errExists)
				}

				if edgeExists {
					log.Tracef("Edge from %s to %s  already exists, doing nothing", *edgeEvent.Event.LsNodeEdge.From, *edgeEvent.Event.LsNodeEdge.To)
					continue
				}

				fromNode, errFrom := repo.CheckNodeExistsById(*edgeEvent.Event.LsNodeEdge.From)
				toNode, errTo := repo.CheckNodeExistsById(*edgeEvent.Event.LsNodeEdge.To)
				if errFrom != nil || errTo != nil {
					log.Fatalf("Error checking if node exists: %s %s", errTo, errFrom)
				}

				if fromNode && toNode {
					log.Infof("Adding edge %s %s", *edgeEvent.Event.LsNodeEdge.From, *edgeEvent.Event.LsNodeEdge.To)
					repo.AddEdge(edge)
					recreateGraph(repo, networkName)
				} else {
					go func(event EdgeEvent) {
						log.Tracef("Adding edge back to queue %s", *event.Event.Key)
						time.Sleep(4 * time.Second)
						edgeChannel <- event
					}(edgeEvent)
				}
			}
		}
	}()
}

func recreateGraph(repo repository.GraphRepository, networkName string) {
	log.Traceln("Recreate graph called")
	if waitMutex.TryLock() {
		go func() {
			defer waitMutex.Unlock()
			time.Sleep(30 * time.Second)
			repo.CreateGraph(networkName)
			log.Infoln("Recreated Neo4j GDS graph")
			repo.CleanupGraph(2)
		}()
	}
}
