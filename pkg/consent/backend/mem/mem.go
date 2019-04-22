// The mem backend package implements the consent storage as well as listener as
// in memory services.
package mem

import (
	//"error"
	"github.com/fission/fission-workflows/pkg/consent"
	"github.com/fission/fission-workflows/pkg/types"
)

// Simple in-memory store of string:ConsentStatus pairs
type Store struct {
	Consent *map[consent.ID]types.ConsentStatus
}

func NewConsentStore() Store {
	return Store{&map[consent.ID]types.ConsentStatus{}}
}

func (cb Store) set(cid consent.ID, status types.ConsentStatus) error {
	//todo check if previous entry exists
	(*cb.Consent)[cid] = status
	return nil
}

func (cb Store) Get(cid consent.ID) types.ConsentStatus {
	val, exist := (*cb.Consent)[cid]

	// if consent not consent status we assume it is granted
	if !exist {
		return types.ConsentStatus{types.ConsentStatus_GRANTED}
	}
	return val
}

func (cb Store) Listen() {
	// TODO
}
