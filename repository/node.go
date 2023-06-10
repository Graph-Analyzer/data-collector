package repository

import (
	"graph-analyzer/data-collector/repository/models"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	log "github.com/sirupsen/logrus"
)

func (r *Graph) AddNode(node models.Node) {
	log.Debugln("Adding node to neo4j")
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	nodes := make([]models.Node, 0, 1)
	nodes = append(nodes, node)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return addNodesFn(tx, &nodes)
	})
	if err != nil {
		log.Fatalf("Error adding node %s", err)
	}
}

func (r *Graph) AddNodes(nodes *[]models.Node) {
	log.Debugln("Adding nodes to neo4j")
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return addNodesFn(tx, nodes)
	})
	if err != nil {
		log.Fatalf("Error adding nodes %s", err)
	}
}

func (r *Graph) CheckNodeExistsById(id string) (bool, error) {
	log.Debugf("Checking if node exists by id %s", id)
	session := r.initReadSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return checkNodeExistsByIdFn(tx, id)
	})
	if err != nil {
		log.Fatalf("Error checking if node exists %s", err)
	}
	return result.(bool), err
}

func (r *Graph) CreateNodeIndex() {
	log.Debugln("Creating (if not exists) neo4j node indexes")
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return createNodeIdIndexFn(tx)
	})

	if err != nil {
		log.Fatalf("Error creating node id index %s", err)
	}

	_, err = session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return createNodeKeyIndexFn(tx)
	})

	if err != nil {
		log.Fatalf("Error creating node key index %s", err)
	}

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return awaitNodeIdIndexCreationFn(tx)
	})

	if err != nil {
		log.Fatalf("Error waiting for node id index creation %s", err)
	}

	_, err = session.ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return awaitNodeKeyIndexCreationFn(tx)
	})

	if err != nil {
		log.Fatalf("Error waiting for node key index creation %s", err)
	}
}

func (r *Graph) DeleteNodeById(id string) {
	log.Debugf("Delete node by id %s", id)
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return deleteNodeByIdFn(tx, id)
	})
	if err != nil {
		log.Fatalf("Error deleting node by id %s", err)
	}
}

func (r *Graph) DeleteNodeByKey(key string) {
	log.Debugf("Delete node by key %s", key)
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return deleteNodeByKeyFn(tx, key)
	})
	if err != nil {
		log.Fatalf("Error deleting node by key %s", err)
	}
}

func (r *Graph) UpdateNodeLabelById(id string, label string) {
	log.Debugf("Updating label of node %s to %s", id, label)
	session := r.initWriteSession()
	defer func(session neo4j.Session) {
		err := session.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(session)

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (any, error) {
		return updateNodeLabelByIdFn(tx, id, label)
	})
	if err != nil {
		log.Fatalf("Error updating node label by id %s", err)
	}
}

func addNodesFn(tx neo4j.Transaction, nodes *[]models.Node) (any, error) {
	log.Tracef("Adding nodes to neo4j %+v", nodes)
	query := `UNWIND $data as props
		CREATE (n:Router)
		SET n = props`
	parameters := map[string]any{
		"data": mapParameters(nodes, mapNode),
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func checkNodeExistsByIdFn(tx neo4j.Transaction, id string) (bool, error) {
	query := `OPTIONAL MATCH (n:Router{RouterID:$routerID})
		RETURN n IS NOT NULL AS Predicate`
	parameters := map[string]any{
		"routerID": id,
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

func createNodeIdIndexFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Creating index_node_router_id index")
	query := `CREATE INDEX index_node_router_id IF NOT EXISTS
		FOR (n:Router)
		ON (n.RouterID)`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func awaitNodeIdIndexCreationFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Awaiting index_node_router_id index")
	query := `CALL db.awaitIndex('index_node_router_id')`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func createNodeKeyIndexFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Creating index_node_router_key index")
	query := `CREATE INDEX index_node_router_key IF NOT EXISTS
		FOR (n:Router)
		ON (n.RouterKey)`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func awaitNodeKeyIndexCreationFn(tx neo4j.Transaction) (any, error) {
	log.Traceln("Awaiting index_node_router_key index")
	query := `CALL db.awaitIndex('index_node_router_key')`
	parameters := map[string]any{}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteNodeByIdFn(tx neo4j.Transaction, id string) (any, error) {
	query := `MATCH (n:Router{RouterId:$routerId})
		DETACH DELETE n`
	parameters := map[string]any{
		"routerId": id,
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func deleteNodeByKeyFn(tx neo4j.Transaction, key string) (any, error) {
	query := `MATCH (n:Router{RouterKey:$routerKey})
		DETACH DELETE n`
	parameters := map[string]any{
		"routerKey": key,
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func updateNodeLabelByIdFn(tx neo4j.Transaction, id string, label string) (any, error) {
	query := `MATCH (n:Router{RouterId:$routerId})
		SET n.Label = $label`
	parameters := map[string]any{
		"routerId": id,
		"label":    label,
	}

	result, err := tx.Run(query, parameters)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func mapNode(node *models.Node) map[string]any {
	return map[string]any{
		"RouterID":  node.RouterID,
		"RouterKey": node.RouterKey,
		"Label":     node.Label,
	}
}
