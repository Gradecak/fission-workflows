package httpclient

import (
	"context"
	"fmt"
	"github.com/gradecak/fission-workflows/pkg/types"
	"net/http"
)

type ConsentAPI struct {
	baseAPI
}

func NewConsentAPI(endpoint string, client http.Client) *ConsentAPI {
	return &ConsentAPI{
		baseAPI: baseAPI{
			endpoint: endpoint,
			client:   client,
		},
	}
}

func (api *ConsentAPI) UpdateFromStatus(ctx context.Context, consentID string, status *types.ConsentStatus) error {
	consentMessage := &types.ConsentMessage{ID: consentID, Status: status}
	return api.Update(ctx, consentMessage)
}

func (api *ConsentAPI) UpdateFromString(ctx context.Context, consentID string, status string) error {
	s, ok := types.ConsentStatus_Status_value[status]
	if !ok {
		return fmt.Errorf("status string %s not valid consent status", status)
	}
	return api.UpdateFromStatus(ctx, consentID, &types.ConsentStatus{types.ConsentStatus_Status(s)})
}

func (api *ConsentAPI) Update(ctx context.Context, msg *types.ConsentMessage) error {
	return callWithJSON(ctx, http.MethodPost, api.formatURL("/consent/update"), msg, nil)
}
