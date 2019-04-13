package consent

import (
	"context"
	"github.com/fission/fission-workflows/pkg/types"
)

type ID = string

type ConsentStore interface {
	Get(ID) types.ConsentStatus
	Listen(context.Context)
	//set should only ever be invoked by the ConsentStore Listener as such
	//it is not exported
	Set(ID, types.ConsentStatus) error
}