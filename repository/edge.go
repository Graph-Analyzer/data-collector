package repository

import (
	"graph-analyzer/data-collector/repository/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
)

func (r *Graph) AddEdge(edge models.Edge) {
	log.Debugln("Adding edge to neo4j")
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	edges := make([]models.Edge, 0, 1)
	edges = append(edges, edge)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return addEdgesFn(tx, &edges)
	})
	if err != nil {
		log.Fatalf("Error adding edge %s", err)
	}
}

func (r *Graph) AddEdges(edges *[]models.Edge) {
	log.Debugln("Adding edges to neo4j")
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return addEdgesFn(tx, edges)
	})
	if err != nil {
		log.Fatalf("Error adding edges %s", err)
	}
}

func (r *Graph) CheckEdgeExistsById(id string) (bool, error) {
	log.Debugf("Checking if edge exists in neo4j %s", id)
	session := r.initReadSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return checkEdgeExistsByIdFn(tx, id)
	})
	if err != nil {
		log.Fatalf("Error checking if edge exists %s", err)
	}
	return result.(bool), err
}

func (r *Graph) CreateEdgeIndex() {
	log.Debugln("Creating (if not exists) neo4j edge indexes")
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return createEdgeConnectsToIdIndexFn(tx)
	})
	if err != nil {
		log.Fatalf("Error creating edge connects to id index %s", err)
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return createEdgeConnectsToKeyIndexFn(tx)
	})
	if err != nil {
		log.Fatalf("Error creating edge connects to key index %s", err)
	}

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return awaitEdgeConnectsToIdIndexCreationFn(tx)
	})
	if err != nil {
		log.Fatalf("Error waiting for edge connects to id index creation %s", err)
	}

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return awaitEdgeConnectsToKeyIndexCreationFn(tx)
	})
	if err != nil {
		log.Fatalf("Error waiting for edge connects to key index creation %s", err)
	}
}

func (r *Graph) DeleteEdgeById(id string) {
	log.Debugf("Delete edge by id %s", id)
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return deleteEdgeByIdFn(tx, id)
	})
	if err != nil {
		log.Fatalf("Error deleting edge by id %s", err)
	}
}

func (r *Graph) DeleteEdgeByKey(key string) {
	log.Debugf("Delete edge by key %s", key)
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return deleteEdgeByKeyFn(tx, key)
	})
	if err != nil {
		log.Fatalf("Error deleting edge by key %s", err)
	}
}

func addEdgesFn(tx neo4j.Transaction, edges *[]models.Edge) (any, error) {
	log.Tracef("Adding edges to neo4j %+v", edges)
	query := `UNWIND $data as props
		MATCH (from:Router {RouterID: props.FromRouterID}),(to:Router {RouterID: props.ToRouterID})
		CALL apoc.create.relationship(
			from,
			props.Label,
			{
				EdgeID: props.EdgeID,
				EdgeKey: props.EdgeKey,
				FromRouterID: props.FromRouterID,
				ToRouterID: props.ToRouterID,
				Weight: props.Weight
			},
			to
		)
		YIELD rel
		RETURN rel`
	parameters := map[string]any{
		"data": mapParameters(edges, mapEdge),
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func checkEdgeExistsByIdFn(tx neo4j.Transaction, id string) (bool, error) {
	query := `OPTIONAL MATCH (:Router)-[r:CONNECTS_TO]->(:Router)
		WHERE r.EdgeID = $edgeID
		RETURN r IS NOT NULL AS Predicate`
	parameters := map[string]any{
		"edgeID": id,
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return false, err
	}
	record, err := result.Single()
	if err != nil {
		return false, err
	}

	return record.Values[0].(bool), nil
}

func createEdgeConnectsToIdIndexFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Creating index_edge_connects_to_id index")
	query := `CREATE INDEX index_edge_connects_to_id IF NOT EXISTS
		FOR ()-[r:CONNECTS_TO]->()
		ON (r.EdgeID)`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func awaitEdgeConnectsToIdIndexCreationFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Awaiting index_edge_connects_to_id index")
	query := `CALL db.awaitIndex('index_edge_connects_to_id')`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func createEdgeConnectsToKeyIndexFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Creating index_edge_connects_to_key index")
	query := `CREATE INDEX index_edge_connects_to_key IF NOT EXISTS
		FOR ()-[r:CONNECTS_TO]->()
		ON (r.EdgeKey)`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func awaitEdgeConnectsToKeyIndexCreationFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Awaiting index_edge_connects_to_key index")
	query := `CALL db.awaitIndex('index_edge_connects_to_key')`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteEdgeByIdFn(tx neo4j.Transaction, id string) (any, error) {
	query := `MATCH (:Router)-[r:CONNECTS_TO]->(:Router) 
		WHERE r.EdgeID=$edgeId 
		DELETE r`
	parameters := map[string]any{
		"edgeId": id,
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteEdgeByKeyFn(tx neo4j.Transaction, key string) (any, error) {
	query := `MATCH (:Router)-[r:CONNECTS_TO]->(:Router) 
		WHERE r.EdgeKey=$edgeKey 
		DELETE r`
	parameters := map[string]any{
		"edgeKey": key,
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func mapEdge(edge *models.Edge) map[string]any {
	return map[string]any{
		"EdgeID":       edge.EdgeID,
		"EdgeKey":      edge.EdgeKey,
		"FromRouterID": edge.FromRouterID,
		"ToRouterID":   edge.ToRouterID,
		"Weight":       edge.Weight,
		"Label":        edge.Label.String(),
	}
}
