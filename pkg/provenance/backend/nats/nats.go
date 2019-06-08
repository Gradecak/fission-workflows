package nats

import (
	fesNATS "github.com/gradecak/fission-workflows/pkg/fes/backend/nats"
	"github.com/gradecak/fission-workflows/pkg/provenance/graph"
	// "github.com/gradecak/fission-workflows/pkg/types"
	"github.com/golang/protobuf/proto"
	"github.com/gradecak/fission-workflows/pkg/util"
	//"github.com/nats-io/go-nats"
	stan "github.com/nats-io/stan.go"
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
	conn, err := stan.Connect(cnf.Cluster, cnf.Client, stan.NatsURL(cnf.URL))
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	return &Publisher{conn}, nil
}

func (p *Publisher) Save(graph *graph.Provenance) error {
	ba, err := proto.Marshal(graph)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	p.Publish(defaultPubPrefix, ba)
	return nil
}
