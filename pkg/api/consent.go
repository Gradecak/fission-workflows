package api

import (
	"context"
	"github.com/fission/fission-workflows/pkg/consent"
	"github.com/fission/fission-workflows/pkg/types"
)

// Consent store consinst of two components, a consentStore in which consent
// information is persisted ready to be queried, and an event store from which
// consent events (REVOKES, PAUSES, etc) are retrieved and processed
type Consent struct {
	cs consent.ConsentStore
	cl consent.ConsentListener
}

func NewConsentAPI(cs consent.ConsentStore, cl consent.ConsentListener) *Consent {
	return &Consent{cs, cl}
}

func (capi *Consent) Query(cid consent.ID) types.ConsentStatus {
	return capi.cs.Get(cid)
}

// Listener starts the Consent Listener and returns. Consent Listener will
// terminate when the provided context invokes the Done() function
func (capi *Consent) Listen(ctx context.Context) {
	go capi.cl.Start(ctx, &capi.cs)
}
