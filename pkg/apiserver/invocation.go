package apiserver

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/gradecak/fission-workflows/pkg/api"
	"github.com/gradecak/fission-workflows/pkg/api/projectors"
	"github.com/gradecak/fission-workflows/pkg/api/store"
	"github.com/gradecak/fission-workflows/pkg/fes"
	"github.com/gradecak/fission-workflows/pkg/fnenv"
	workflowFnenv "github.com/gradecak/fission-workflows/pkg/fnenv/workflows"
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/gradecak/fission-workflows/pkg/types/validate"
	"github.com/gradecak/fission-workflows/pkg/util"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"sort"
	"time"
)

var (
	invocationTime = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "apiserver",
		Subsystem: "invocation",
		Name:      "time",
		Help:      "Duration of an invocation from start to a finished state.",
		Objectives: map[float64]float64{
			0.25: 0.0001,
			0.5:  0.0001,
			0.9:  0.0001,
			1:    0.0001,
		},
	})
	invocationConcurrent = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "apiserver",
		Subsystem: "invocation",
		Name:      "concurrent",
		Help:      "Number of active controllers",
	})
)

func init() {
	prometheus.MustRegister(invocationTime, invocationConcurrent)
}

// Invocation is responsible for all functionality related to managing invocations.
type Invocation struct {
	api         *api.Invocation
	invocations *store.Invocations
	workflows   *store.Workflows
	fnenv       *workflowFnenv.Runtime
	backend     fes.Backend
}

func NewInvocation(api *api.Invocation, invocations *store.Invocations, workflows *store.Workflows, backend fes.Backend) WorkflowInvocationAPIServer {
	return &Invocation{
		api:         api,
		invocations: invocations,
		workflows:   workflows,
		fnenv:       workflowFnenv.NewRuntime(api, invocations, workflows),
		backend:     backend,
	}
}

func (gi *Invocation) Validate(ctx context.Context, spec *types.WorkflowInvocationSpec) (*empty.Empty, error) {
	err := validate.WorkflowInvocationSpec(spec)
	if err != nil {
		return nil, toErrorStatus(err)
	}
	return &empty.Empty{}, nil
}

func (gi *Invocation) Invoke(ctx context.Context, spec *types.WorkflowInvocationSpec) (*types.ObjectMetadata, error) {
	// TODO go through same runtime as InvokeSync
	// Check if the workflow required by the invocation exists
	wf, err := gi.workflows.GetWorkflow(spec.GetWorkflowId())
	if err != nil {
		return nil, err
	}
	spec.Workflow = wf

	eventID, err := gi.api.Invoke(spec, api.WithContext(ctx))
	if err != nil {
		return nil, toErrorStatus(err)
	}

	return &types.ObjectMetadata{Id: eventID}, nil
}

func (gi *Invocation) InvokeSync(ctx context.Context, spec *types.WorkflowInvocationSpec) (*types.WorkflowInvocation, error) {
	invocationConcurrent.Inc()
	start := time.Now()
	defer func() {
		invocationTime.Observe(float64(time.Since(start)))
		invocationConcurrent.Dec()
	}()
	wfi, err := gi.fnenv.InvokeWorkflow(spec, fnenv.WithContext(ctx))
	if err != nil {
		return nil, toErrorStatus(err)
	}
	return wfi, nil
}

func (gi *Invocation) Cancel(ctx context.Context, objectMetadata *types.ObjectMetadata) (*empty.Empty, error) {
	err := gi.api.Cancel(objectMetadata.GetId())
	if err != nil {
		return nil, toErrorStatus(err)
	}

	return &empty.Empty{}, nil
}

func (gi *Invocation) Get(ctx context.Context, objectMetadata *types.ObjectMetadata) (*types.WorkflowInvocation, error) {
	wi, err := gi.invocations.GetInvocation(objectMetadata.GetId())
	if err != nil {
		return nil, toErrorStatus(err)
	}
	return wi, nil
}

func (gi *Invocation) List(ctx context.Context, query *InvocationListQuery) (*WorkflowInvocationList, error) {
	var invocations []string
	as := gi.invocations.List()
	for _, aggregate := range as {
		if aggregate.Type != types.TypeInvocation {
			logrus.Errorf("Invalid type in invocation invocations: %v", aggregate.Format())
			continue
		}

		if len(query.Workflows) > 0 {
			// TODO make more efficient (by moving list queries to invocations)
			entity, err := gi.invocations.GetAggregate(aggregate)
			if err != nil {
				logrus.Errorf("List: failed to fetch %v from invocations: %v", aggregate, err)
				continue
			}
			wfi := entity.(*types.WorkflowInvocation)
			if !contains(query.Workflows, wfi.GetSpec().GetWorkflowId()) {
				continue
			}
		}

		invocations = append(invocations, aggregate.Id)
	}
	return &WorkflowInvocationList{Invocations: invocations}, nil
}

func (gi *Invocation) AddTask(ctx context.Context, req *AddTaskRequest) (*empty.Empty, error) {
	invocation, err := gi.invocations.GetInvocation(req.GetInvocationID())
	if err != nil {
		return nil, toErrorStatus(err)
	}
	if err := gi.api.AddTask(invocation.ID(), req.Task); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (gi *Invocation) Events(ctx context.Context, md *types.ObjectMetadata) (*ObjectEvents, error) {
	events, err := gi.backend.Get(projectors.NewWorkflowAggregate(md.Id))
	if err != nil {
		return nil, toErrorStatus(err)
	}

	// TODO this should not be this cumbersome
	wi, err := gi.invocations.GetInvocation(md.Id)
	if err != nil {
		return nil, toErrorStatus(err)
	}

	// Fold task events into invocation events
	for _, task := range wi.GetStatus().GetTasks() {
		taskEvents, err := gi.taskEvents(task.ID())
		if err != nil {
			return nil, toErrorStatus(fmt.Errorf("failed to fetch task events: %v", err))
		}
		events = append(events, taskEvents...)
	}

	sort.SliceStable(events, func(i, j int) bool {
		return util.CmpProtoTimestamps(events[i].GetTimestamp(), events[j].GetTimestamp())
	})

	return &ObjectEvents{
		Metadata: md,
		Events:   events,
	}, nil
}

func (gi *Invocation) taskEvents(taskRunID string) ([]*fes.Event, error) {
	return gi.backend.Get(projectors.NewTaskRunAggregate(taskRunID))
}

func contains(haystack []string, needle string) bool {
	for i := 0; i < len(haystack); i++ {
		if haystack[i] == needle {
			return true
		}
	}
	return false
}
