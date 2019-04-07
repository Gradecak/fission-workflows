package consent

import (
	"github.com/fission/fission-workflows/pkg/types"
)

type ConsentBackend struct {
	Consent map[ID]types.ConsentStatus
}

func (cb *ConsentBackend) Set(cid ID, status types.ConsentStatus) error {
	//todo check if previous entry exists
	cb.Consent[cid] = status
	return nil
}

func (cb *ConsentBackend) Get(cid ID) types.ConsentStatus {
	val, exist := cb.Consent[cid]

	// if consent not consent status we assume it is granted
	if !exist {
		return types.ConsentStatus{types.ConsentStatus_GRANTED}
	}

	return val
}
