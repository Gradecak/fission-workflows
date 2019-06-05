package filestore

import (
	// "github.com/fission/fission-workflows/pkg/provenance"
	"github.com/fission/fission-workflows/pkg/provenance/graph"
	// "github.com/fission/fission-workflows/pkg/types"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"os"
)

type Publisher struct {
	*os.File
}

func NewPublisher(fName string) (*Publisher, error) {
	file, err := os.Create(fName)
	if err != nil {
		return nil, err
	}
	return &Publisher{file}, nil
}

func (p Publisher) Save(g *graph.Provenance) error {
	// err := provenance.ValidateGraph(g, graph.Node_UNDEF)
	// if err != nil {
	// 	logrus.Error(err.Error())
	// 	return err
	// }
	ba, err := proto.Marshal(g)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	_, err = p.Write(ba)
	return err
}
