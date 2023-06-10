package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalGexf(t *testing.T) {
	xmlData := []byte(`
		<?xml version="1.0" encoding="UTF-8"?>
		<gexf xmlns="http://www.gexf.net/1.2draft" version="1.2">
			<graph>
				<nodes>
					<node id="1" label="Node 1"/>
					<node id="2" label="Node 2"/>
				</nodes>
				<edges>
					<edge id="1" source="1" target="2"/>
				</edges>
			</graph>
		</gexf>
	`)

	graph, err := UnmarshalGexf(xmlData)

	assert.NoError(t, err)

	assert.NotNil(t, graph.Graph)
	assert.NotNil(t, graph.Graph.Nodes)
	assert.NotNil(t, graph.Graph.Nodes.Nodes)
	assert.NotNil(t, graph.Graph.Edges)
	assert.NotNil(t, graph.Graph.Edges.Edges)

	assert.Equal(t, 2, len(graph.Graph.Nodes.Nodes))
	assert.Equal(t, 1, len(graph.Graph.Edges.Edges))
	assert.Equal(t, "Node 1", graph.Graph.Nodes.Nodes[0].Label)
	assert.Equal(t, "Node 2", graph.Graph.Nodes.Nodes[1].Label)
}
