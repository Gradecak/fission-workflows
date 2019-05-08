package fnenv

import (
	"fmt"
	"github.com/fission/fission-workflows/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	//	"strings"
	//	"errors"
	"sync"
	"time"
)

const (
	defaultTimeout = time.Duration(1) * time.Minute
)

var (
	fnResolved = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "fnenv",
		Subsystem: "fission",
		Name:      "functions_resolved_total",
		Help:      "Total number of Fission functions resolved",
	}, []string{"fnenv"})
)

func init() {
	prometheus.MustRegister(fnResolved)
}

// MetaResolver contacts function execution runtime clients to resolve the function definitions to concrete function ids.
//
// ParseTask definitions (See types/TaskDef) can contain the following function reference:
// - `<name>` : the function is currently resolved to one of the clients
// - `<client>:<name>` : forces the client that the function needs to be resolved to.
//
// Future:
// - Instead of resolving just to one client, resolve function for all clients, and apply a priority or policy
//   for scheduling (overhead vs. load)
//
type MetaResolver struct {
	clients map[string]RuntimeResolver
	timeout time.Duration
}

func NewMetaResolver(client map[string]RuntimeResolver) *MetaResolver {
	return &MetaResolver{
		clients: client,
		timeout: defaultTimeout,
	}
}

func (ps *MetaResolver) Resolve(targetFn string) (types.FnRef, error) {
	ref, err := types.ParseFnRef(targetFn)
	if err != nil {
		return types.FnRef{}, err
	}

	if ref.Runtime != "" {
		return ps.resolveForRuntime(ref.Runtime, ref)
	}

	waitFor := len(ps.clients)
	resolved := make(chan types.FnRef, waitFor)
	defer close(resolved)
	wg := sync.WaitGroup{}
	wg.Add(waitFor)
	for cName := range ps.clients {
		go func(cName string) {
			def, err := ps.resolveForRuntime(cName, ref)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"err":     err,
					"runtime": cName,
					"fn":      targetFn,
				}).Debug("Function not found.")
			} else {
				resolved <- def
			}
			wg.Done()
		}(cName)
	}
	wg.Wait() // for all clients to resolve

	// For now just select the first resolved
	select {
	case result := <-resolved:
		return result, nil
	default:
		return types.FnRef{}, fmt.Errorf("failed to resolve function '%s' using clients '%v'", targetFn, ps.clients)
	}
}

func (ps *MetaResolver) resolveForRuntime(runtime string, ref types.FnRef) (types.FnRef, error) {
	dst, ok := ps.clients[runtime]
	if !ok {
		return types.FnRef{}, ErrInvalidRuntime
	}
	rsv, err := dst.Resolve(ref)
	if err != nil {
		return types.FnRef{}, err
	}

	fnResolved.WithLabelValues(runtime).Inc()

	//set the zone by parsing the -XX suffix from the function ID
	zone := ref.GenZone()
	return types.FnRef{
		Runtime:   runtime,
		Namespace: ref.Namespace,
		ID:        rsv,
		Zone:      zone,
	}, nil
}

//
// Helper functions
//

// ResolveTask resolved all the tasks in the provided workflow spec.
//
// In case there are no functions resolved, an empty slice is returned.
func ResolveTasks(ps Resolver, tasks map[string]*types.TaskSpec) (map[string]*types.FnRef, error) {
	var flattened []*types.TaskSpec
	for _, v := range tasks {
		flattened = append(flattened, v)
	}
	return ResolveTask(ps, flattened...)
}

// ResolveTask resolved the interpreted workflow from a given spec.
//
// It returns a map consisting of the original functionRef as the key.
func ResolveTask(ps Resolver, tasks ...*types.TaskSpec) (map[string]*types.FnRef, error) {
	// Check for duplicates
	uniqueTasks := map[string]*types.TaskSpec{}
	for _, t := range tasks {
		if _, ok := uniqueTasks[t.FunctionRef]; !ok {
			uniqueTasks[t.FunctionRef] = t
		}
	}

	logrus.Debugf("UNIQUE TASKS %+v", uniqueTasks)

	var lastErr error
	wg := sync.WaitGroup{}
	wg.Add(len(uniqueTasks))
	resolved := map[string]*types.FnRef{}
	resolvedC := make(chan []sourceFnRef, len(uniqueTasks))

	// ResolveTask each task in the workflow definition in parallel
	for k, t := range uniqueTasks {
		go func(k string, t *types.TaskSpec, tc chan []sourceFnRef) {
			err := resolveTask(ps, k, t, tc)
			if err != nil {
				lastErr = err
			}
			wg.Done()
		}(k, t, resolvedC)
	}

	// Close channel when all tasks have resolved
	go func() {
		wg.Wait()
		close(resolvedC)
	}()

	// Store results of the resolved tasks
	for r := range resolvedC {
		for _, t := range r {
			resolved[t.src] = t.FnRef
		}
	}

	// A little bit of a hack... in the case of multizone tasks, in case the
	// base function ref (eg, hello) doesnt resolve, but a zone variant does
	// (eg hello-nl) we must change the FunctionRef of the taskSpec to the zone
	// variant in order to be able to execute the Task
	for _, t := range tasks {
		if uTask, ok := uniqueTasks[t.FunctionRef]; ok {
			t.FunctionRef = uTask.FunctionRef
		}
	}

	if lastErr != nil {
		return nil, lastErr
	}
	return resolved, nil
}

// resolveTaskAndInputs traverses the inputs of a task to resolve nested functions and workflows.
func resolveTask(ps Resolver, id string, task *types.TaskSpec, resolvedC chan []sourceFnRef) error {
	if task == nil || resolvedC == nil {
		return nil
	}

	toResolve := []string{task.FunctionRef}
	resolved := []sourceFnRef{}

	if constr := task.GetExecConstraints(); constr != nil {
		if constr.GetMultiZone() {
			zoneVariants := types.GenZoneVariants(task.FunctionRef)
			toResolve = append(toResolve, zoneVariants...)
		}
	}

	// resolve all of the functions across all zones for the given task
	for i, r := range toResolve {
		t, err := ps.Resolve(r)
		if err != nil {
			logrus.Errorf("Error resolving fn %v, reason: %+v --- Skipping", r, err)

			// if the current function ref couldnt be resolved, move
			// the TaskSpec.FnRef on to the next function in the
			// list. useful when its a multizone task and the base
			// task function does not exist eg hello-world -->
			// doesn't exist .. hello-world-nl --> exists
			if r == task.FunctionRef && i < len(toResolve)-1 {
				task.FunctionRef = toResolve[i+1]
			}
			continue
		}
		resolved = append(resolved, sourceFnRef{
			src:   r,
			FnRef: &t,
		})
	}

	if len(resolved) < 1 {
		return fmt.Errorf("Could not resolve fnRef: %v", task.FunctionRef)
	}

	resolvedC <- resolved
	return nil
}

// sourceFnRef wraps the FnRef with the source field which contains the original function target.
type sourceFnRef struct {
	*types.FnRef
	src string
}
