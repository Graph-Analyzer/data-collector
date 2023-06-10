package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
)

func (r *Graph) CreateGraph(networkName string) {
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return createGraphFn(tx, networkName)
	})
	if err != nil {
		log.Fatalf("Error adding graph %s", err)
	}
}

func (r *Graph) CleanupGraph(skip int64) {
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return cleanupGraphFn(tx, skip)
	})
	if err != nil {
		log.Fatalf("Error cleanup graph %s", err)
	}
}

// Create undirected unweighted projection of graph
// Analyzed data should not consider direction and weights
func createGraphFn(tx neo4j.Transaction, networkName string) (any, error) {
	query := `CALL gds.graph.project(
			$graphName,
			'Router',
			{
				TYPE: {
					type: 'CONNECTS_TO',
					orientation: 'UNDIRECTED',
					aggregation: 'SINGLE'
				}
			}
	)`
	parameters := map[string]any{
		"graphName": fmt.Sprintf("%s:%s", networkName, uuid.New().String()),
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func cleanupGraphFn(tx neo4j.Transaction, skip int64) (any, error) {
	query := `CALL gds.graph.list()
		YIELD graphName, creationTime
		RETURN graphName, creationTime
		ORDER BY creationTime DESC SKIP $skip`
	parameters := map[string]any{
		"skip": skip,
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		query := `CALL gds.graph.drop($graphName)
			YIELD graphName;`
		parameters := map[string]any{
			"graphName": result.Record().Values[0],
		}

		_, err := tx.Run(query, parameters)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
