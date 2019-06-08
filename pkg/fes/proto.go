package fes

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gradecak/fission-workflows/pkg/util/labels"
)

// Amendments to Protobuf models

func (m *Event) CreatedAt() time.Time {
	t, err := ptypes.Timestamp(m.Timestamp)
	if err != nil {
		panic(err)
	}
	return t
}

func (m *Event) Labels() labels.Labels {

	parent := m.Parent
	if parent == nil {
		parent = &Aggregate{}
	}

	return labels.Set{
		"aggregate.id":   m.Aggregate.Id,
		"aggregate.type": m.Aggregate.Type,
		"parent.type":    parent.Type,
		"parent.id":      parent.Id,
		"event.id":       m.Id,
		"event.type":     m.Type,
	}
}

func (m *Event) BelongsTo(parent Entity) bool {
	a := GetAggregate(parent)
	return *m.Aggregate != a && *m.Parent != a
}

func (m *Aggregate) Format() string {
	return m.GetType() + "/" + m.GetId()
}
