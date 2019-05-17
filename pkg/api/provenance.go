package api

import (
	"github.com/fission/fission-workflows/pkg/provenance"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/sirupsen/logrus"
)

type Provenance struct {
	provenance.Store
}

func NewProvenance(store provenance.Store) *Provenance {
	return &Provenance{store}
}

// TODO Currently the provenance graph is regenerated for every finished
// invocation, however, the only difference is the root node. As an optimisation
// we could cache workflow provenance DAGs with the workflow ID
func (p Provenance) GenerateProvenance(id string, spec *types.WorkflowSpec) error {
	graph := &types.Node{Type: types.Node_ROOT, Tag: id}

	edge := &types.Edge{EdgeType: types.DEFAULT_EDGE_TYPE, Dst: spec.ToProvenance()}
	graph.Edges = append(graph.Edges, edge)

	// validate the graph structure before notifying listeners
	err := provenance.ValidateGraph(graph, types.Node_UNDEF)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	// publish the graph to listeners
	p.Save(graph)
	return nil
}
