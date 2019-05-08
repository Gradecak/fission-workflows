package nats

import (
	"github.com/fission/fission-workflows/pkg/consent"
	//import for the nats configs to avoid duplicate delarations
	fesNATs "github.com/fission/fission-workflows/pkg/fes/backend/nats"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/fission/fission-workflows/pkg/util"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/sirupsen/logrus"
)

const (
	defaultClient    = "consent"
	defaultCluster   = "consent-cluster"
	defaultURL       = nats.DefaultURL
	defaultSubPrefix = "CONSENT"
)

type Store struct {
	Consent *map[consent.ID]*types.ConsentStatus
	// subHandler func(m *stan.Msg)
	nats stan.Conn
	sub  stan.Subscription
}

func NewNatsConsentStore(cnf fesNATs.Config) (*Store, error) {
	logrus.Debugf("CONSENT nats config %+v", cnf)

	cnf.Client = util.UID()
	if cnf.URL == "" {
		cnf.URL = defaultURL
	}
	if cnf.Cluster == "" {
		cnf.Cluster = defaultCluster
	}
	conn, err := fesNATs.ConnectNats(cnf)

	if err != nil {
		return nil, err
	}

	store := &Store{
		Consent: &map[consent.ID]*types.ConsentStatus{},
		nats:    conn,
	}

	return store, nil
}

func (cb Store) Get(cid consent.ID) types.ConsentStatus {
	val, exist := (*cb.Consent)[cid]

	// if consent not consent status we assume it is granted
	if !exist {
		return types.ConsentStatus{types.ConsentStatus_GRANTED}
	}
	return *val
}

func (cb *Store) Listen() {
	sub, err := cb.nats.Subscribe(defaultSubPrefix,
		func(m *stan.Msg) {
			logrus.Info("Consent Event Recieved..")
			msg, err := toConsentMsg(m.Data)
			if err != nil {
				logrus.Error("Failed to unmarshal Consent Message")
			}
			cb.setConsent(msg)
		},
		stan.DeliverAllAvailable())

	if err != nil {
		logrus.Errorf("Failed to subscribe to %v channel", defaultSubPrefix)
	}
	cb.sub = sub
	//sub.Unsubscribe()
}

func (cb Store) setConsent(m *consent.ConsentMessage) {
	(*cb.Consent)[m.GetId()] = m.GetStatus()
}

func toConsentMsg(data []byte) (*consent.ConsentMessage, error) {
	cmsg := &consent.ConsentMessage{}
	err := proto.Unmarshal(data, cmsg)
	if err != nil {
		return nil, err
	}
	return cmsg, nil
}
