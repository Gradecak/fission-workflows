package provenance

import (
	"errors"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/sirupsen/logrus"
)

type Store interface {
	Save(*types.Node) error
}

func ValidateGraph(n *types.Node, prev types.Node_Type) error {
	switch n.Type {
	case types.Node_ROOT:

		if prev != types.Node_UNDEF {
			return errors.New("Root node contains predecessor")
		}

		edges := n.GetEdges()
		var err error = nil
		for _, e := range edges {
			err = ValidateGraph(e.Dst, types.Node_ROOT)
			if err != nil {
				return err
			}
		}
	case types.Node_DATAFLOW:
		if prev != types.Node_ROOT {
			return errors.New("Dataflow node predecessor not root")
		}

	}
	return nil
}
