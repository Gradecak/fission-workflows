package consent

import (
	"github.com/fission/fission-workflows/pkg/types"
)

type ID = string

type ConsentStore interface {
	Get(ID) *types.ConsentStatus
	Listen()
	Set(*types.ConsentMessage) error
}
