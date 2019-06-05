package provenance

import (
	// "errors"
	"github.com/fission/fission-workflows/pkg/provenance/graph"
	// "github.com/fission/fission-workflows/pkg/types"
	// "github.com/sirupsen/logrus"
)

type Store interface {
	Save(*graph.Provenance) error
}

func ValidateGraph(n *graph.Node, prev graph.Node_Type) error {
	// TODO
	return nil
}
