package monitor

import (
	// "github.com/gradecak/fission-workflows/pkg/controller/ctrl"
	"github.com/prometheus/client_golang/prometheus"
	// "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type InvocationMonitor struct {
	activeMu          sync.RWMutex
	activeInvocations map[string]time.Time
	activeCount       int
	queuedMu          sync.RWMutex
	queuedCount       int
	queuedInvocations map[string]time.Time
}

var (
	monitorActiveCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "invocation",
		Subsystem: "monitor",
		Name:      "active",
		Help:      "Number of running invocations",
	})
	monitorScheduledCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "invocation",
		Subsystem: "monitor",
		Name:      "scheduled",
		Help:      "Number of scheduled",
	})
	InvocationTime = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "invocation",
		Subsystem: "monitor",
		Name:      "time",
		Help:      "Duration of an invocation in various states.",
		Objectives: map[float64]float64{
			0:    0.0001,
			0.01: 0.0001,
			0.1:  0.0001,
			0.25: 0.0001,
			0.5:  0.0001,
			0.75: 0.0001,
			0.9:  0.0001,
			1:    0.0001,
		},
	}, []string{"type"})
)

func init() {
	prometheus.MustRegister(monitorScheduledCount, monitorActiveCount, InvocationTime)
}

func NewInvocationMonitor() *InvocationMonitor {
	return &InvocationMonitor{
		activeMu:          sync.RWMutex{},
		activeInvocations: make(map[string]time.Time),
		activeCount:       0,
		queuedMu:          sync.RWMutex{},
		queuedInvocations: make(map[string]time.Time),
		queuedCount:       0,
	}
}

func (m *InvocationMonitor) AddQueued(invID string) {
	m.queuedMu.Lock()
	if _, ok := m.queuedInvocations[invID]; !ok {
		m.queuedInvocations[invID] = time.Now()
		m.queuedCount++
		monitorScheduledCount.Set(float64(m.queuedCount))
	}
	m.queuedMu.Unlock()
}

func (m *InvocationMonitor) RunQueued(invID string) {
	m.queuedMu.Lock()
	if t, ok := m.queuedInvocations[invID]; ok {
		delete(m.queuedInvocations, invID)
		m.queuedCount--

		m.activeMu.Lock()
		m.activeInvocations[invID] = t
		m.activeCount++
		m.activeMu.Unlock()

		defer func() {
			monitorScheduledCount.Set(float64(m.queuedCount))
			InvocationTime.WithLabelValues("queued").Observe(float64(time.Since(t)))
			monitorActiveCount.Set(float64(m.activeCount))
		}()
	}
	m.queuedMu.Unlock()

}

func (m *InvocationMonitor) Remove(invID string) {
	// if in queuedQueue
	m.queuedMu.Lock()
	if t, ok := m.queuedInvocations[invID]; ok {
		delete(m.queuedInvocations, invID)
		m.queuedCount--
		monitorScheduledCount.Set(float64(m.queuedCount))
		InvocationTime.WithLabelValues("queued").Observe(float64(time.Since(t)))
		m.queuedMu.Unlock()
		return
	}
	m.queuedMu.Unlock()

	// if in activeQueue
	m.activeMu.Lock()
	if t, ok := m.activeInvocations[invID]; ok {
		delete(m.activeInvocations, invID)
		m.activeCount--
		monitorActiveCount.Set(float64(m.activeCount))
		InvocationTime.WithLabelValues("running").Observe(float64(time.Since(t)))
		m.activeMu.Unlock()
		return
	}
	m.activeMu.Unlock()

}

func (m *InvocationMonitor) ActiveCount() int {
	return m.activeCount
}

func (m *InvocationMonitor) QueuedCount() int {
	return m.queuedCount
}

// func (m *InvocationMonitor) HasBacklog() bool {
// 	// return float64(m.activeCount/ctrl.MAX_PARALLEL_CONTROLLERS) >= 0.85 && m.queuedCount > 0
// 	return float64()
// }
