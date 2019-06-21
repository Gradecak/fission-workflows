package controller

import (
	"context"
	// "fmt"
	"github.com/gradecak/fission-workflows/pkg/api"
	"github.com/gradecak/fission-workflows/pkg/api/store"
	"github.com/gradecak/fission-workflows/pkg/controller/ctrl"
	"github.com/gradecak/fission-workflows/pkg/controller/monitor"
	// "github.com/gradecak/fission-workflows/pkg/fes"
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	POLL_INTERVAL = time.Millisecond * 500
)

type AdmissionController struct {
	runOnce *sync.Once
	sensors []ctrl.Sensor
}

func NewAdmissionController(invocations *store.Invocations, api *api.Invocation, maxParallel int) *AdmissionController {
	return &AdmissionController{
		runOnce: &sync.Once{},
		sensors: []ctrl.Sensor{
			NewQueuedInvocationSensor(invocations, api, api.GetMonitor(), maxParallel),
			NewQueuedInvocationPollSensor(invocations, api, POLL_INTERVAL, api.GetMonitor(), maxParallel),
		},
	}
}

func (c *AdmissionController) Run() {
	c.runOnce.Do(func() {
		go c.run()
	})
}

func (c *AdmissionController) run() {
	for _, sensor := range c.sensors {
		err := sensor.Start(nil)
		if err != nil {
			panic(err)
		}
	}
}

type QueuedInvocationSensor struct {
	invocations *store.Invocations
	maxParallel int
	api         *api.Invocation
	monitor     *monitor.InvocationMonitor
	done        func()
	closeC      <-chan struct{}
}

func NewQueuedInvocationSensor(invocations *store.Invocations, api *api.Invocation, monitor *monitor.InvocationMonitor, maxParallel int) *QueuedInvocationSensor {
	ctx, done := context.WithCancel(context.Background())
	return &QueuedInvocationSensor{
		invocations: invocations,
		maxParallel: maxParallel,
		api:         api,
		monitor:     monitor,
		done:        done,
		closeC:      ctx.Done(),
	}
}

func (s *QueuedInvocationSensor) Start(eq ctrl.EvalQueue) error {
	go s.Run(eq)
	return nil
}

func (s *QueuedInvocationSensor) Run(ctrl.EvalQueue) {
	sub := s.invocations.GetQueuedInvocationUpdates()
	if sub == nil {
		logrus.Warn("Invocation store does not support pubsub.")
		return
	}
	for {
		select {
		case msg := <-sub.Ch:
			notification, err := sub.ToNotification(msg)
			if err != nil {
				logrus.Warnf("Failed to convert pubsub message to notification: %v", err)
			}
			// if there is invocations waiting to be queued defer to
			// the poll sensor to schedule them first as then we are
			// guaranteed fairness in scheduling
			if float64(s.monitor.ActiveCount()/s.maxParallel) >= 0.85 && s.monitor.QueuedCount() > 0 {
				if s.monitor.ActiveCount() < s.maxParallel {
					s.api.Schedule(notification.Aggregate.GetId())
				}
			}
		case <-s.closeC:
			err := sub.Close()
			if err != nil {
				logrus.Error(err)
			}
			logrus.Info("Notification listener stopped.")
			return
		}
	}
}

func (s *QueuedInvocationSensor) Close() error {
	s.done()
	return nil
}

type QueuedInvocationPollSensor struct {
	*ctrl.PollSensor
	invocations *store.Invocations
	maxParallel int
	monitor     *monitor.InvocationMonitor
	api         *api.Invocation
}

func NewQueuedInvocationPollSensor(invocations *store.Invocations, api *api.Invocation,
	interval time.Duration, monitor *monitor.InvocationMonitor, maxParallel int) *QueuedInvocationPollSensor {

	sens := &QueuedInvocationPollSensor{
		invocations: invocations,
		maxParallel: maxParallel,
		monitor:     monitor,
		api:         api,
	}
	sens.PollSensor = ctrl.NewPollSensor(interval, sens.Poll)
	return sens
}

func (s *QueuedInvocationPollSensor) Poll(ctrl.EvalQueue) {
	for _, aggregate := range s.invocations.List() {
		if aggregate.Type != types.TypeInvocation {
			continue
		}
		wf, err := s.invocations.GetInvocation(aggregate.GetId())
		if err != nil {
			logrus.Warnf("Could not retrieve entity from invocations store: %v", aggregate)
			continue
		}
		// only concerned with invocaitons that are in IN_PROGRESS state
		if !wf.GetStatus().Queued() {
			continue
		}

		// if controller not permitted wait until next polling period
		if s.monitor.ActiveCount() >= s.maxParallel {
			return
		}
		s.api.Schedule(aggregate.Id)
	}
}
