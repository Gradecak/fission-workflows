package api

import (
	"github.com/fission/fission-workflows/pkg/provenance"
	"github.com/fission/fission-workflows/pkg/provenance/graph"
	"github.com/fission/fission-workflows/pkg/types"
	// "github.com/sirupsen/logrus"
)

type Provenance struct {
	provenance.Store
}

func NewProvenance(store provenance.Store) *Provenance {
	return &Provenance{store}
}

// TODO Currently the provenance graph is regenerated for every finished
// invocation, however, the sometimes only difference is the root node. As an optimisation
// we could cache workflow provenance DAGs with the workflow ID
func (p Provenance) GenerateProvenance(wfi *types.WorkflowInvocation) error {
	prov := graph.GenProvenance(wfi)
	p.Save(prov)
	return nil
}
