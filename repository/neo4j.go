package repository

import (
	"graph-analyzer/data-collector/repository/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
)

type GraphRepository interface {
	AddEdge(edge models.Edge)
	AddEdges(edges *[]models.Edge)
	AddNode(node models.Node)
	AddNodes(nodes *[]models.Node)
	CheckEdgeExistsById(id string) (bool, error)
	CheckNodeExistsById(id string) (bool, error)
	CleanupGraph(skip int64)
	CreateGraph(networkName string)
	CreateEdgeIndex()
	CreateNodeIndex()
	DeleteAll()
	DeleteEdgeById(id string)
	DeleteEdgeByKey(key string)
	DeleteNodeById(id string)
	DeleteNodeByKey(key string)
	UpdateNodeLabelById(id string, label string)
}

type Graph struct {
	driver neo4j.Driver
}

func NewGraphRepository(driver neo4j.Driver) GraphRepository {
	repo := &Graph{
		driver: driver,
	}
	repo.DeleteAll()
	repo.CleanupGraph(0)

	return repo
}

func (r *Graph) DeleteAll() {
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		if result, err := deleteAllRelationshipsFn(tx); err != nil {
			return result, err
		}

		return deleteAllNodesFn(tx)
	})
	if err != nil {
		log.Fatalf("Error deleting all node and edges %s", err)
	}
}

func deleteAllRelationshipsFn(tx neo4j.Transaction) (any, error) {
	query := `CALL apoc.periodic.iterate(
		'MATCH ()-[r]->() RETURN id(r) AS id', 
		'MATCH ()-[r]->() WHERE id(r)=id DELETE r', 
		{batchSize: 5000})`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteAllNodesFn(tx neo4j.Transaction) (any, error) {
	query := `CALL apoc.periodic.iterate(
		'MATCH (n) RETURN id(n) AS id', 
		'MATCH (n) WHERE id(n)=id DELETE n', 
		{batchSize: 5000})`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Graph) initReadSession() neo4j.Session {
	session := r.driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeRead,
	})

	return session
}

func (r *Graph) initWriteSession() neo4j.Session {
	session := r.driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})

	return session
}

func mapParameters[P any](
	parameters *[]P,
	mappingFn func(*P) map[string]any,
) []map[string]any {
	var result = make([]map[string]any, len(*parameters))

	for i, parameter := range *parameters {
		result[i] = mappingFn(&parameter)
	}

	return result
}
