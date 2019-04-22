package api

import (
	"github.com/fission/fission-workflows/pkg/provenance"
	"github.com/fission/fission-workflows/pkg/types"
)

type Provenance struct {
	provenance.Store
}

func NewProvenance(store provenance.Store) *Provenance {
	return &Provenance{store}
}

func (p Provenance) GenerateProvenance(id string, spec *types.WorkflowSpec) error {
	graph := &types.Node{Type: types.Node_ROOT, Tag: id}
	graph.Edges = append(graph.Edges, spec.ToProvenance())
	//TODO validate graph
	p.Save(graph)
	return nil
}
