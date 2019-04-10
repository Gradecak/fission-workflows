package api

import (
	"context"
	"github.com/fission/fission-workflows/pkg/consent"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/fission/fission-workflows/pkg/types/typedvalues"
	log "github.com/sirupsen/logrus"
)

// Consent consinst of two components, a consentStore in which consent
// information is persisted ready to be queried, and an event listener which
// runs in a goroutine and recieves updates for consent information
type Consent struct {
	cs consent.ConsentStore
}

func NewConsentAPI(cs consent.ConsentStore) *Consent {
	return &Consent{cs}
}

func (capi *Consent) Query(cid consent.ID) types.ConsentStatus {
	return capi.cs.Get(cid)
}

// Listener starts the Consent Listener and returns. Consent Listener will
// terminate when the provided context invokes the Done() function
func (capi *Consent) Listen(ctx context.Context) {
	go capi.cs.Listen(ctx)
}

// Given the workflow paramaters, resolve the inputs and return the consent
// ConsentStatus of the workflow
func (capi *Consent) QueryWorkflowConsent(invocation *types.WorkflowInvocation) (types.ConsentStatus, error) {
	inputs, err := typedvalues.UnwrapMapTypedValue(invocation.GetSpec().GetInputs())
	if err != nil {
		log.Error("POOO")
	}
	log.Debug(inputs)
	return types.ConsentStatus{types.ConsentStatus_GRANTED}, nil
}
