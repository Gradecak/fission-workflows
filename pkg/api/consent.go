package api

import (
	"github.com/fission/fission-workflows/pkg/consent"
	"github.com/fission/fission-workflows/pkg/types"
	//"github.com/sirupsen/logrus"
)

// Consent consinst of two components, a consentStore in which consent
// information is persisted ready to be queried, and an event listener which
// runs in a goroutine and recieves updates for consent information
type Consent struct {
	consent.ConsentStore
}

func NewConsentAPI(cs consent.ConsentStore) *Consent {
	return &Consent{cs}
}

func (capi *Consent) Query(cid consent.ID) types.ConsentStatus {
	return capi.Get(cid)
}

// Listener starts the Consent Listener and returns. Consent Listener will
// terminate when the provided context invokes the Done() function
func (capi *Consent) WatchConsent() {
	capi.Listen()
}

// Given the workflow paramaters, resolve the inputs and return the consent
// ConsentStatus of the workflow
func (capi *Consent) QueryWorkflowConsent(spec *types.WorkflowInvocationSpec) types.ConsentStatus {
	consentId := spec.GetConsentId()
	status := capi.Get(consentId)
	return status
}
