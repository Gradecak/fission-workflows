package scheduler

import (
	"fmt"
	math "math"
	"time"

	"math/rand"
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/gradecak/fission-workflows/pkg/types/graph"
)

// HorizonPolicy is the default policy of the workflow engine. It solely schedules tasks that are on the scheduling horizon.
//
// The scheduling horizon is the set of tasks that only depend on tasks that have already completed.
// If a task has failed this policy simply fails the workflow
type HorizonPolicy struct {
}

func NewHorizonPolicy() *HorizonPolicy {
	return &HorizonPolicy{}
}

func (p *HorizonPolicy) Evaluate(invocation *types.WorkflowInvocation) (*Schedule, error) {
	schedule := &Schedule{InvocationId: invocation.ID(), CreatedAt: ptypes.TimestampNow()}

	// If there are failed tasks halt the workflow
	if failedTasks := getFailedTasks(invocation); len(failedTasks) > 0 {
		for _, failedTask := range failedTasks {
			msg := fmt.Sprintf("Task '%v' failed", failedTask.ID())
			if err := failedTask.GetStatus().GetError(); err != nil {
				msg = err.Message
			}
			schedule.Abort = newAbortAction(msg)
		}
		return schedule, nil
	}

	// Find and schedule all tasks on the scheduling horizon
	openTasks := getOpenTasks(invocation)
	depGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	horizon := graph.Roots(depGraph)
	for _, node := range horizon {
		schedule.AddRunTask(newRunTaskAction(node.(*graph.TaskInvocationNode).Task().ID()))
	}
	return schedule, nil
}

// PrewarmAllPolicy is the policy with the most aggressive form of prewarming.
//
// The policy, like the HorizonPolicy, schedules all tasks on the scheduling horizon optimistically.
// Similarly, it also fails workflow invocations immediately if a task has failed
//
// However, on top of the HorizonPolicy, this policy prewarms tasks aggressively. Any unstarted task not on the
// scheduling horizon will be prewarmed.
//
// This policy does not try to infer runtimes or cold starts; instead, it prewarms with a static duration.
type PrewarmAllPolicy struct {
	coldStartDuration time.Duration
}

func NewPrewarmAllPolicy(coldstartDuration time.Duration) *PrewarmAllPolicy {
	return &PrewarmAllPolicy{coldStartDuration: coldstartDuration}
}

func (p *PrewarmAllPolicy) Evaluate(invocation *types.WorkflowInvocation) (*Schedule, error) {
	schedule := &Schedule{InvocationId: invocation.ID(), CreatedAt: ptypes.TimestampNow()}

	// If there are failed tasks halt the workflow
	if failedTasks := getFailedTasks(invocation); len(failedTasks) > 0 {
		for _, failedTask := range failedTasks {
			msg := fmt.Sprintf("Task '%v' failed", failedTask.ID())
			if err := failedTask.GetStatus().GetError(); err != nil {
				msg = err.Message
			}
			schedule.Abort = newAbortAction(msg)
		}
		return schedule, nil
	}

	// Find and schedule all tasks on the scheduling horizon
	openTasks := getOpenTasks(invocation)
	depGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	horizon := graph.Roots(depGraph)
	for _, node := range horizon {
		taskRun := node.(*graph.TaskInvocationNode)
		schedule.AddRunTask(newRunTaskAction(taskRun.TaskInvocation.ID()))
		delete(openTasks, taskRun.GetMetadata().GetId())
	}

	// Prewarm all other tasks
	expectedAt := time.Now().Add(p.coldStartDuration)
	for _, task := range openTasks {
		schedule.AddPrepareTask(newPrepareTaskAction(task.ID(), expectedAt))
	}
	return schedule, nil
}

// PrewarmHorizonPolicy is the policy with the most aggressive form of prewarming.
//
// The policy, like the HorizonPolicy, schedules all tasks on the scheduling horizon optimistically.
// Similarly, it also fails workflow invocations immediately if a task has failed
//
// However, on top of the HorizonPolicy, tries to policy prewarms tasks aggressively. Any unstarted task on the
// prewarm horizon will be prewarmed.
//
// This policy does not try to infer runtimes or cold starts; instead, it prewarms with a static duration.
type PrewarmHorizonPolicy struct {
	coldStartDuration time.Duration
}

func NewPrewarmHorizonPolicy(coldstartDuration time.Duration) *PrewarmHorizonPolicy {
	return &PrewarmHorizonPolicy{coldStartDuration: coldstartDuration}
}

func (p *PrewarmHorizonPolicy) Evaluate(invocation *types.WorkflowInvocation) (*Schedule, error) {
	schedule := &Schedule{InvocationId: invocation.ID(), CreatedAt: ptypes.TimestampNow()}

	// If there are failed tasks halt the workflow
	if failedTasks := getFailedTasks(invocation); len(failedTasks) > 0 {
		for _, failedTask := range failedTasks {
			msg := fmt.Sprintf("Task '%v' failed", failedTask.ID())
			if err := failedTask.GetStatus().GetError(); err != nil {
				msg = err.Message
			}
			schedule.Abort = newAbortAction(msg)
		}
		return schedule, nil
	}

	// Find and schedule all tasks on the scheduling horizon
	openTasks := getOpenTasks(invocation)
	depGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	horizon := graph.Roots(depGraph)
	for _, node := range horizon {
		taskRun := node.(*graph.TaskInvocationNode)
		schedule.AddRunTask(newRunTaskAction(taskRun.TaskInvocation.ID()))
		delete(openTasks, taskRun.GetMetadata().GetId())
	}

	// Prewarm all tasks on the prewarm horizon
	// Note: we are mutating openTasks!
	expectedAt := time.Now().Add(p.coldStartDuration)
	prewarmDepGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	prewarmHorizon := graph.Roots(prewarmDepGraph)
	for _, node := range prewarmHorizon {
		taskRun := node.(*graph.TaskInvocationNode)
		schedule.AddPrepareTask(newPrepareTaskAction(taskRun.Task().ID(), expectedAt))
	}

	return schedule, nil
}

// Horizon Multizone Policy is similar to Horizon Policy, however, it take into
// account the Zone hints specified in the Workflow Specificiation
type HorizonMultiZonePolicy struct {
	random *rand.Rand
}

func NewHorizonMZPolicy() *HorizonMultiZonePolicy {
	seed := rand.NewSource(time.Now().Unix())
	return &HorizonMultiZonePolicy{rand.New(seed)}
}

func (p *HorizonMultiZonePolicy) Evaluate(invocation *types.WorkflowInvocation) (*Schedule, error) {
	schedule := &Schedule{InvocationId: invocation.ID(), CreatedAt: ptypes.TimestampNow()}

	// If there are failed tasks halt the workflow
	if failedTasks := getFailedTasks(invocation); len(failedTasks) > 0 {
		for _, failedTask := range failedTasks {
			msg := fmt.Sprintf("Task '%v' failed", failedTask.ID())
			if err := failedTask.GetStatus().GetError(); err != nil {
				msg = err.Message
			}
			schedule.Abort = newAbortAction(msg)
		}
		return schedule, nil
	}

	// Find and schedule all tasks on the scheduling horizon
	openTasks := getOpenTasks(invocation)
	depGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	horizon := graph.Roots(depGraph)

	for _, node := range horizon {
		task := node.(*graph.TaskInvocationNode).Task()
		taskAction := newRunTaskAction(task.ID())

		// set the preferred function execution environment
		var target *types.FnRef
		if task.GetSpec().GetExecConstraints().GetMultiZone() {
			if rf, ok := task.GetZoneLock(); ok {
				// if zone lock is specified in spec use that
				target = rf
			} else if ref, ok := invocation.GetPreferredZone(task); ok {
				// Run time zone hints take precedence over hone hints
				// provided in spec
				target = ref
			} else if ref, ok := task.GetZoneHint(); ok {
				// fall back on workflow spec specified zone hint
				target = ref
			} else {
				//fall back on random fnref
				refs := task.GetAltFnRefs()
				target = refs[p.random.Intn(len(refs))]
			}

		}

		// update count for that particular fnRef
		taskAction.Pref = target
		schedule.AddRunTask(taskAction)
	}

	return schedule, nil
}

type MzHorizonLRUWarmPolicy struct {
	random *rand.Rand
	envsMu *sync.Mutex
	envs   map[*types.FnRef]time.Time
	// timeout before an environment is no longer consindered "warm"
	warmTimeout time.Duration
	// a static 'estimate' for task execution
	avgExecTime time.Duration
}

func NewMzHorizonLRUWarmPolicy() *MzHorizonLRUWarmPolicy {
	seed := rand.NewSource(time.Now().Unix())
	return &MzHorizonLRUWarmPolicy{
		random:      rand.New(seed),
		envsMu:      &sync.Mutex{},
		envs:        make(map[*types.FnRef]time.Time),
		warmTimeout: time.Minute * 2,
		avgExecTime: time.Second * 1,
	}
}

// type FnRank struct {
// 	ref  *types.FnRef
// 	rank float64
// }

func (p *MzHorizonLRUWarmPolicy) Evaluate(invocation *types.WorkflowInvocation) (*Schedule, error) {
	schedule := &Schedule{InvocationId: invocation.ID(), CreatedAt: ptypes.TimestampNow()}

	// If there are failed tasks halt the workflow
	if failedTasks := getFailedTasks(invocation); len(failedTasks) > 0 {
		for _, failedTask := range failedTasks {
			msg := fmt.Sprintf("Task '%v' failed", failedTask.ID())
			if err := failedTask.GetStatus().GetError(); err != nil {
				msg = err.Message
			}
			schedule.Abort = newAbortAction(msg)
		}
		return schedule, nil
	}

	// Find and schedule all tasks on the scheduling horizon
	openTasks := getOpenTasks(invocation)
	depGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	horizon := graph.Roots(depGraph)

	for _, node := range horizon {
		task := node.(*graph.TaskInvocationNode).Task()
		taskAction := newRunTaskAction(task.ID())
		p.envsMu.Lock()
		// set the preferred function execution environment
		var rankedRef = struct {
			fn   *types.FnRef
			rank float64
		}{
			fn:   nil,
			rank: math.MaxFloat64,
		}
		if task.GetSpec().GetExecConstraints().GetMultiZone() {
			if rf, ok := task.GetZoneLock(); ok {
				// if zone lock is specified in spec use that
				rankedRef.fn = rf
			} else if ref, ok := invocation.GetPreferredZone(task); ok {
				// Run time zone hints take precedence over hone hints
				// provided in spec
				rankedRef.fn = ref
			} else if ref, ok := task.GetZoneHint(); ok {
				// fall back on workflow spec specified zone hint
				rankedRef.fn = ref
			} else {
				refs := task.GetAltFnRefs()
				for _, r := range refs {
					if lastUsed, ok := p.envs[r]; ok {
						var rank float64 = (0.5 * float64(p.warmTimeout/time.Since(lastUsed))) +
							(0.5 * float64(time.Since(lastUsed)/p.avgExecTime))
						if rank < rankedRef.rank {
							rankedRef.rank = rank
							rankedRef.fn = r
						}
					}
				}
			}
		}
		// update count for that particular fnRef
		p.envs[rankedRef.fn] = time.Now()
		p.envsMu.Unlock()
		taskAction.Pref = rankedRef.fn //target
		schedule.AddRunTask(taskAction)
	}

	return schedule, nil
}

type MzHorizonLRUPolicy struct {
	random *rand.Rand
	envsMu *sync.Mutex
	envs   map[string]*types.FnRef
}

func NewMzHorizonLRUPolicy() *MzHorizonLRUPolicy {
	seed := rand.NewSource(time.Now().Unix())
	return &MzHorizonLRUPolicy{
		random: rand.New(seed),
		envsMu: &sync.Mutex{},
		envs:   make(map[string]*types.FnRef),
	}
}

func (p *MzHorizonLRUPolicy) Evaluate(invocation *types.WorkflowInvocation) (*Schedule, error) {
	schedule := &Schedule{InvocationId: invocation.ID(), CreatedAt: ptypes.TimestampNow()}

	// If there are failed tasks halt the workflow
	if failedTasks := getFailedTasks(invocation); len(failedTasks) > 0 {
		for _, failedTask := range failedTasks {
			msg := fmt.Sprintf("Task '%v' failed", failedTask.ID())
			if err := failedTask.GetStatus().GetError(); err != nil {
				msg = err.Message
			}
			schedule.Abort = newAbortAction(msg)
		}
		return schedule, nil
	}

	// Find and schedule all tasks on the scheduling horizon
	openTasks := getOpenTasks(invocation)
	depGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	horizon := graph.Roots(depGraph)

	for _, node := range horizon {
		task := node.(*graph.TaskInvocationNode).Task()
		taskAction := newRunTaskAction(task.ID())
		// set the preferred function execution environment
		var target *types.FnRef
		if task.GetSpec().GetExecConstraints().GetMultiZone() {
			if rf, ok := task.GetZoneLock(); ok {
				// if zone lock is specified in spec use that
				target = rf
			} else if ref, ok := invocation.GetPreferredZone(task); ok {
				// Run time zone hints take precedence over hone hints
				// provided in spec
				target = ref
			} else if ref, ok := task.GetZoneHint(); ok {
				// fall back on workflow spec specified zone hint
				target = ref
			} else {
				if fnref, ok := p.envs[task.Status.GetBaseFnName()]; ok {
					target = fnref
				} else {
					target = task.GetAltFnRefs()[0]
				}
			}
		}
		// update count for that particular fnRef
		p.envsMu.Lock()
		p.envs[task.Status.GetBaseFnName()] = target
		p.envsMu.Unlock()
		taskAction.Pref = target
		schedule.AddRunTask(taskAction)
	}

	return schedule, nil
}

type MzHorizonRRPolicy struct {
	random *rand.Rand
	envsMu *sync.Mutex
	envs   map[*types.FnRef][]*types.FnRef
}

func NewMzHorizonRRPolicy() *MzHorizonRRPolicy {
	seed := rand.NewSource(time.Now().Unix())
	return &MzHorizonRRPolicy{
		random: rand.New(seed),
		envsMu: &sync.Mutex{},
		envs:   make(map[*types.FnRef][]*types.FnRef),
	}
}

func (p *MzHorizonRRPolicy) Evaluate(invocation *types.WorkflowInvocation) (*Schedule, error) {
	schedule := &Schedule{InvocationId: invocation.ID(), CreatedAt: ptypes.TimestampNow()}

	// If there are failed tasks halt the workflow
	if failedTasks := getFailedTasks(invocation); len(failedTasks) > 0 {
		for _, failedTask := range failedTasks {
			msg := fmt.Sprintf("Task '%v' failed", failedTask.ID())
			if err := failedTask.GetStatus().GetError(); err != nil {
				msg = err.Message
			}
			schedule.Abort = newAbortAction(msg)
		}
		return schedule, nil
	}

	// Find and schedule all tasks on the scheduling horizon
	openTasks := getOpenTasks(invocation)
	depGraph := graph.Parse(graph.NewTaskInstanceIterator(openTasks))
	horizon := graph.Roots(depGraph)

	for _, node := range horizon {
		task := node.(*graph.TaskInvocationNode).Task()
		taskAction := newRunTaskAction(task.ID())

		// set the preferred function execution environment
		var target *types.FnRef
		if task.GetSpec().GetExecConstraints().GetMultiZone() {
			if rf, ok := task.GetZoneLock(); ok {
				// if zone lock is specified in spec use that
				target = rf
			} else if ref, ok := invocation.GetPreferredZone(task); ok {
				// Run time zone hints take precedence over hone hints
				// provided in spec
				target = ref
			} else if ref, ok := task.GetZoneHint(); ok {
				// fall back on workflow spec specified zone hint
				target = ref
			} else {
				p.envsMu.Lock()
				if baseRef := task.GetStatus().GetFnRef(); baseRef != nil {
					envs, ok := p.envs[baseRef]
					if !ok {
						p.envs[baseRef] = task.GetAltFnRefs()
						envs = p.envs[baseRef]
					}
					// round robin selection
					target, rest := envs[0], envs[1:]
					p.envs[baseRef] = append(rest, target)
				}
				p.envsMu.Unlock()
			}

		}

		// update count for that particular fnRef
		taskAction.Pref = target
		schedule.AddRunTask(taskAction)
	}

	return schedule, nil
}

// ************************************************************
//                      END OF POLICIES
// ************************************************************

func getFailedTasks(invocation *types.WorkflowInvocation) []*types.TaskInvocation {
	var failedTasks []*types.TaskInvocation
	for _, task := range invocation.TaskInvocations() {
		if task.GetStatus().GetStatus() == types.TaskInvocationStatus_FAILED {
			failedTasks = append(failedTasks, task)
		}
	}
	return failedTasks
}

func getOpenTasks(invocation *types.WorkflowInvocation) map[string]*types.TaskInvocation {
	openTasks := map[string]*types.TaskInvocation{}
	for id, task := range invocation.Tasks() {
		taskRun, ok := invocation.TaskInvocation(id)
		if !ok {
			taskRun = &types.TaskInvocation{
				Metadata: types.NewObjectMetadata(id),
				Spec:     types.NewTaskInvocationSpec(invocation, task, time.Now()),
				Status: &types.TaskInvocationStatus{
					Status: types.TaskInvocationStatus_UNKNOWN,
				},
			}
		}
		if taskRun.GetStatus().GetStatus() == types.TaskInvocationStatus_UNKNOWN {
			openTasks[id] = taskRun
		}
	}
	return openTasks
}
