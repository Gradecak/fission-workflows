package consent

import (
	"context"
	"github.com/fission/fission-workflows/pkg/types"
)

type ID = string

type ConsentStore interface {
	Get(ID) types.ConsentStatus
	Set(ID, types.ConsentStatus) error
}

// Listener polls for ConsentStatus messages from a given source and on recieved
// events updates the consent storage.
type ConsentListener interface {
	Start(context.Context, *ConsentStore) error
}
