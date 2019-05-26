package apiserver

import (
	"errors"
	"github.com/fission/fission-workflows/pkg/api"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Consent struct {
	api *api.Consent
}

func NewConsent(api *api.Consent) ConsentAPIServer {
	return &Consent{api}
}

func (c *Consent) Update(ctx context.Context, cm *types.ConsentMessage) (*empty.Empty, error) {
	if cm.GetID() == "" {
		return nil, errors.New("Id recieved is empty")
	}
	if cm.GetStatus() == nil {
		return nil, errors.New("Recieved status is empty")
	}
	logrus.Infof("Got message %v", cm)
	oldStatus := c.api.Query(cm.GetID())
	c.api.Set(cm)
	logrus.Infof("Status changed for ID: %v ~~~ %v --> %v", cm.GetID(), oldStatus, c.api.Query(cm.GetID()))
	return &empty.Empty{}, nil
}
