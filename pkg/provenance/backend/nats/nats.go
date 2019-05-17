package nats

import (
	fesNATS "github.com/fission/fission-workflows/pkg/fes/backend/nats"
	//"github.com/fission/fission-workflows/pkg/provenance"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/fission/fission-workflows/pkg/util"
	"github.com/golang/protobuf/proto"
	//"github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/sirupsen/logrus"
)

const (
	defaultPubPrefix = "PROVENANCE"
)

type Publisher struct {
	stan.Conn
}

func NewPublisher(cnf fesNATS.Config) (*Publisher, error) {
	cnf.Client = util.UID()
	conn, err := fesNATS.ConnectNats(cnf)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return &Publisher{conn}, nil
}

func (p *Publisher) Save(graph *types.Node) error {
	ba, err := proto.Marshal(graph)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	p.Publish(defaultPubPrefix, ba)
	return nil
}
