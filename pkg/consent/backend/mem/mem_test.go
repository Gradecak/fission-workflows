package mem

import (
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var ids = []string{"bob", "alice", "tom"}

func TestInsertRetrieve(t *testing.T) {
	store := NewConsentStore()
	for _, id := range ids {
		s := &types.ConsentStatus{types.ConsentStatus_REVOKED}
		store.Set(&types.ConsentMessage{
			ID:     id,
			Status: s,
		})
		status := store.Get(id)
		assert.Equal(t, status, &types.ConsentStatus{types.ConsentStatus_REVOKED})
	}
	assert.Equal(t, len(ids), len(*store.Consent))
}
