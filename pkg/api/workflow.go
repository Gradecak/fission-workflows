package api

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/ptypes"
	"github.com/gradecak/fission-workflows/pkg/api/events"
	"github.com/gradecak/fission-workflows/pkg/api/projectors"
	"github.com/gradecak/fission-workflows/pkg/fes"
	"github.com/gradecak/fission-workflows/pkg/fnenv"
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/gradecak/fission-workflows/pkg/types/validate"
	"github.com/gradecak/fission-workflows/pkg/util"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

// Workflow contains the API functionality for controlling workflow definitions.
// This includes creating and parsing workflows.
type Workflow struct {
	es       fes.Backend
	resolver fnenv.Resolver
}

// NewWorkflowAPI creates the Workflow API.
func NewWorkflowAPI(esClient fes.Backend, resolver fnenv.Resolver) *Workflow {
	return &Workflow{esClient, resolver}
}

// Create creates a new workflow based on the provided workflowSpec.
// The function either returns the id of the workflow or an error.
// The error can be a validate.Err, proto marshall error, or a fes error.
// TODO check if id already exists
func (wa *Workflow) Create(workflow *types.WorkflowSpec, opts ...CallOption) (string, error) {
	cfg := parseCallOptions(opts)
	err := validate.WorkflowSpec(workflow)
	if err != nil {
		return "", err
	}

	// If no id is provided generate an id
	id := workflow.ForceId
	if len(id) == 0 {
		id = fmt.Sprintf("wf-%s", util.UID())
	}

	// Ensure that if provenance metadata is present describing the
	// predecessor dependancy, the predecssor is a valid
	// Workflow/Cloudfunction in the system
	if predecessor := workflow.GetPredecessor(); predecessor != "" {
		logrus.Debug("Predecessor is not nil")
		_, err := wa.resolver.Resolve(predecessor)
		if err != nil {
			logrus.Debug("Cannot resolve the function")
			return "", fmt.Errorf("Cannot resolve the predecessor target")
		}
	}

	event, err := fes.NewEvent(projectors.NewWorkflowAggregate(id), &events.WorkflowCreated{
		Spec: workflow,
	})
	if err != nil {
		return "", err
	}

	// If part of a span, add trace metadata to the event.
	span := opentracing.SpanFromContext(cfg.ctx)
	if span != nil {
		err = opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap,
			opentracing.TextMapCarrier(event.Metadata))
		if err != nil {
			logrus.Warnf("Failed to inject tracer context into event: %v", err)
		}
	}

	err = wa.es.Append(event)
	if err != nil {
		return "", err
	}

	return id, nil
}

// Delete marks a workflow as deleted, making it unavailable to any future interactions.
// This also means that subsequent invocations for this workflow will fail.
// If the API fails to append the event to the event store, it will return an error.
func (wa *Workflow) Delete(workflowID string) error {
	if len(workflowID) == 0 {
		return validate.NewError("workflowID", errors.New("id should not be empty"))
	}

	event, err := fes.NewEvent(projectors.NewWorkflowAggregate(workflowID), &events.WorkflowDeleted{})
	if err != nil {
		return err
	}
	event.Hints = &fes.EventHints{Completed: true}
	return wa.es.Append(event)
}

// Parse processes the workflow to resolve any ambiguity.
// Currently, this means that all the function references are resolved to function identifiers. For convenience
// this function returns the new WorkflowStatus. If the API fails to append the event to the event store,
// it will return an error.
func (wa *Workflow) Parse(workflow *types.Workflow) (map[string]*types.TaskStatus, error) {
	if err := validate.WorkflowSpec(workflow.Spec); err != nil {
		return nil, err
	}

	resolvedFns, err := fnenv.ResolveTasks(wa.resolver, workflow.Spec.Tasks)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve tasks in workflow: %v", err)
	}

	taskStatuses := map[string]*types.TaskStatus{}
	for id, t := range workflow.Spec.Tasks {
		var ref *types.FnRef
		// Grab the FnRef for the task run target, if it failed to
		// resolve, check if a zone variant was resolved and use that
		// instead
		if fnRef, ok := resolvedFns[t.FunctionRef]; ok {
			ref = fnRef
		} else {
			zoneVariants := types.GenZoneVariants(t.FunctionRef)
			for _, zoneVar := range zoneVariants {
				if fnRef, ok := resolvedFns[zoneVar]; ok {
					ref = fnRef
					break
				}
			}
			if ref == nil {
				return nil, fmt.Errorf("Could not resolve run target for Task %v", t.FunctionRef)
			}
		}
		taskStatuses[id] = &types.TaskStatus{
			UpdatedAt:  ptypes.TimestampNow(),
			FnRef:      ref,
			Status:     types.TaskStatus_READY,
			AltFnRefs:  resolvedFns,
			BaseFnName: t.FunctionRef,
		}
	}

	event, err := fes.NewEvent(projectors.NewWorkflowAggregate(workflow.ID()), &events.WorkflowParsed{
		Tasks: taskStatuses,
	})
	if err != nil {
		return nil, err
	}
	err = wa.es.Append(event)
	if err != nil {
		return nil, err
	}

	return taskStatuses, nil
}
