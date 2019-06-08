package provenance

import (
	// "errors"
	"github.com/gradecak/fission-workflows/pkg/provenance/graph"
	// "github.com/gradecak/fission-workflows/pkg/types"
	// "github.com/sirupsen/logrus"
)

type Store interface {
	Save(*graph.Provenance) error
}

func ValidateGraph(n *graph.Node, prev graph.Node_Type) error {
	// TODO
	return nil
}
