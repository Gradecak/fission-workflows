package types

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gradecak/fission-workflows/pkg/types/typedvalues"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// Types other than specified in protobuf
const (
	InputMain    = "default"
	InputBody    = "body"
	InputHeaders = "headers"
	InputQuery   = "query"
	InputMethod  = "method"
	InputParent  = "_parent"

	typedValueShortMaxLen = 32
	WorkflowAPIVersion    = "v1"

	TypeWorkflow   = "workflow"
	TypeInvocation = "invocation"
	TypeTaskRun    = "taskrun"

	EvictAfter = time.Minute
)

// InvocationEvent
var invocationFinalStates = []WorkflowInvocationStatus_Status{
	WorkflowInvocationStatus_ABORTED,
	WorkflowInvocationStatus_SUCCEEDED,
	WorkflowInvocationStatus_FAILED,
}

var taskFinalStates = []TaskInvocationStatus_Status{
	TaskInvocationStatus_FAILED,
	TaskInvocationStatus_ABORTED,
	TaskInvocationStatus_SKIPPED,
	TaskInvocationStatus_SUCCEEDED,
}

//
// Error
//

func (m *Error) Error() string {
	return m.Message
}

//
// WorkflowInvocation
//

func (m *WorkflowInvocation) ID() string {
	return m.GetMetadata().GetId()
}

func (m *WorkflowInvocation) Copy() *WorkflowInvocation {
	return proto.Clone(m).(*WorkflowInvocation)
}

func (m *WorkflowInvocation) Type() string {
	return TypeInvocation
}

func (m *WorkflowInvocation) LastUpdated() (time.Time, error) {
	return m.GetStatus().LastUpdated()
}

func (m WorkflowInvocation) Evictable() bool {
	ts, err := m.LastUpdated()
	if err != nil {
		return false // cannot determine if we should evict
	}
	return time.Since(ts) >= EvictAfter
}

func (m *WorkflowInvocation) Workflow() *Workflow {
	return m.GetSpec().GetWorkflow()
}

func (m *WorkflowInvocation) GetDataflowSpec() *DataFlowSpec {
	return m.GetSpec().GetWorkflow().GetSpec().GetDataflow()
}

func (m *WorkflowInvocation) GetWorkflowSpec() *WorkflowSpec {
	return m.GetSpec().GetWorkflow().GetSpec()
}

// TODO how do we know which tasks are not being run
func (m *WorkflowInvocation) TaskInvocation(id string) (*TaskInvocation, bool) {
	ti, ok := m.Status.Tasks[id]
	return ti, ok
}

func (m *WorkflowInvocation) TaskInvocations() map[string]*TaskInvocation {
	tasks := map[string]*TaskInvocation{}
	for id := range m.Status.Tasks {
		task, _ := m.TaskInvocation(id)
		tasks[id] = task
	}
	return tasks
}

func (m *WorkflowInvocation) Task(id string) (*Task, bool) {
	if dtasks := m.GetStatus().GetDynamicTasks(); dtasks != nil {
		dtask, ok := dtasks[id]
		if ok {
			return dtask, true
		}
	}
	return m.Workflow().Task(id)
}

// Tasks gets all tasks in a workflow. This includes the dynamic tasks added during
// the invocation.
func (m *WorkflowInvocation) Tasks() map[string]*Task {
	tasks := map[string]*Task{}
	if m == nil {
		return tasks
	}
	if wf := m.Workflow(); wf != nil {
		for _, task := range m.Workflow().Tasks() {
			tasks[task.ID()] = task
		}
	}
	for _, task := range m.GetStatus().GetDynamicTasks() {
		tasks[task.ID()] = task
	}
	return tasks
}

func (m *WorkflowInvocation) GetPreferredZone(t *Task) (*FnRef, bool) {
	taskName := t.Metadata.GetId()
	// check if we got runtime hints for task zone
	if hints := m.GetSpec().GetTaskHints(); hints != nil {
		zone, ok := hints[taskName]
		if !ok {
			return nil, false
		}
		ref, err := t.GetZoneVariant(zone)
		if err != nil {
			return nil, false
		}
		return ref, true
	}
	return nil, false
}

func (m *WorkflowInvocation) HasConsentId() bool {
	return len(m.GetSpec().GetConsentId()) > 0
}

type WorkflowInvocations []*WorkflowInvocation

func (s WorkflowInvocations) Len() int      { return len(s) }
func (s WorkflowInvocations) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type ByQueueTime struct{ WorkflowInvocations }

func (s ByQueueTime) Less(i, j int) bool {
	x, _ := ptypes.Timestamp(s.WorkflowInvocations[i].GetMetadata().GetCreatedAt())
	y, _ := ptypes.Timestamp(s.WorkflowInvocations[j].GetMetadata().GetCreatedAt())
	return y.After(x)
}

//
// WorkflowInvocationStatus
//

func (m *WorkflowInvocationStatus) ToTaskStatus() *TaskInvocationStatus {
	var statusMapping = map[WorkflowInvocationStatus_Status]TaskInvocationStatus_Status{
		WorkflowInvocationStatus_UNKNOWN:     TaskInvocationStatus_UNKNOWN,
		WorkflowInvocationStatus_SCHEDULED:   TaskInvocationStatus_SCHEDULED,
		WorkflowInvocationStatus_IN_PROGRESS: TaskInvocationStatus_IN_PROGRESS,
		WorkflowInvocationStatus_SUCCEEDED:   TaskInvocationStatus_SUCCEEDED,
		WorkflowInvocationStatus_FAILED:      TaskInvocationStatus_FAILED,
		WorkflowInvocationStatus_ABORTED:     TaskInvocationStatus_ABORTED,
	}

	return &TaskInvocationStatus{
		Status:        statusMapping[m.Status],
		Error:         m.Error,
		UpdatedAt:     m.UpdatedAt,
		Output:        m.Output,
		OutputHeaders: m.OutputHeaders,
	}
}

// Finished returns true if the invocation is in a terminal state.
func (m WorkflowInvocationStatus) Finished() bool {
	for _, event := range invocationFinalStates {
		if event == m.GetStatus() {
			return true
		}
	}
	return false
}

func (m WorkflowInvocationStatus) Successful() bool {
	return m.GetStatus() == WorkflowInvocationStatus_SUCCEEDED
}

func (m WorkflowInvocationStatus) Running() bool {
	return m.GetStatus() == WorkflowInvocationStatus_IN_PROGRESS
}
func (m WorkflowInvocationStatus) Queued() bool {
	return m.GetStatus() == WorkflowInvocationStatus_SCHEDULED
}

func (m WorkflowInvocationStatus) LastUpdated() (time.Time, error) {
	return ptypes.Timestamp(m.GetUpdatedAt())
}

//
// TaskInvocation
//

func (m *TaskInvocation) ID() string {
	return m.GetMetadata().GetId()
}

func (m *TaskInvocation) Copy() *TaskInvocation {
	return proto.Clone(m).(*TaskInvocation)
}

func (m *TaskInvocation) Type() string {
	return TypeTaskRun
}

func (m *TaskInvocation) Task() *Task {
	return m.GetSpec().GetTask()
}

//
// TaskInvocationStatus
//

func (ti TaskInvocationStatus) Finished() bool {
	for _, event := range taskFinalStates {
		if event == ti.Status {
			return true
		}
	}
	return false
}

func (ti TaskInvocationStatus) Successful() bool {
	return ti.GetStatus() == TaskInvocationStatus_SUCCEEDED
}

//
// Task
//

func (m *Task) ID() string {
	return m.GetMetadata().GetId()
}

func (m *Task) GetZoneLock() (*FnRef, bool) {
	z := m.GetSpec().GetExecConstraints().GetZoneLock()
	if z != Zone_UNDEF {
		if ref, err := m.GetZoneVariant(z); err != nil {
			return ref, true
		}
	}
	return nil, false
}

func (m *Task) GetZoneHint() (*FnRef, bool) {
	z := m.GetSpec().GetExecConstraints().GetZoneHint()
	if z != Zone_UNDEF {
		if ref, err := m.GetZoneVariant(z); err != nil {
			return ref, true
		}
	}
	return nil, false
}

func (m *Task) GetAltFnRefs() []*FnRef {
	r := []*FnRef{}
	for _, ref := range m.GetStatus().GetAltFnRefs() {
		r = append(r, ref)
	}
	return r
}

func (m *Task) GetZoneVariant(z Zone) (*FnRef, error) {
	if z == Zone_UNDEF {
		return nil, fmt.Errorf("Zone provided is undefined")
	}
	altFnRefs := m.GetStatus().GetAltFnRefs()
	ref, ok := altFnRefs[ZoneVariant(m.GetStatus().BaseFnName, z)]
	if !ok {
		return nil, fmt.Errorf("Could not find zone %s in available environemnts", Zone_name[int32(z)])
	}
	return ref, nil
}

//
// TaskSpec
//

func (m *TaskSpec) Input(key string, val *typedvalues.TypedValue) *TaskSpec {
	if len(m.Inputs) == 0 {
		m.Inputs = map[string]*typedvalues.TypedValue{}
	}
	m.Inputs[key] = val

	return m
}

func (m *TaskSpec) Parent() (string, bool) {
	var parent string
	var present bool
	for id, params := range m.Requires {
		if params.Type == TaskDependencyParameters_DYNAMIC_OUTPUT {
			present = true
			parent = id
			break
		}
	}
	return parent, present
}

func (m *TaskSpec) Require(taskID string, opts ...*TaskDependencyParameters) *TaskSpec {
	if m.Requires == nil {
		m.Requires = map[string]*TaskDependencyParameters{}
	}
	var params *TaskDependencyParameters
	if len(opts) > 0 {
		params = opts[0]
	}

	m.Requires[taskID] = params
	return m
}

//
//func (m *TaskSpec) Overlay(overlay *TaskSpec) *TaskSpec {
//	nt := proto.Clone(m).(*TaskSpec)
//	nt.Await = overlay.Await
//	nt.Requires = overlay.Requires
//	return nt
//}

//
// Workflow
//
func (m *Workflow) ID() string {
	return m.GetMetadata().GetId()
}

func (m *Workflow) Copy() *Workflow {
	return proto.Clone(m).(*Workflow)
}

func (m *Workflow) Type() string {
	return TypeWorkflow
}

// Note: this only retrieves the statically, top-level defined tasks
// TODO just store entire task in status
func (m *Workflow) Task(id string) (*Task, bool) {
	//var ok bool
	//spec, ok := m.Spec.Tasks[id]
	//if !ok {
	//	return nil, false
	//}
	//var task *TaskStatus
	//if m.Status.Tasks != nil {
	//	task, ok = m.Status.Tasks[id]
	//}
	//if !ok {
	//	task = &TaskStatus{
	//		UpdatedAt: ptypes.TimestampNow(),
	//	}
	//}
	//
	//return &Task{
	//	Metadata: &ObjectMetadata{
	//		Id:        id,
	//		CreatedAt: m.Metadata.CreatedAt,
	//	},
	//	Spec:   spec,
	//	Status: task,
	//}, true

	// First the status to get the latest task state
	ts := m.GetStatus().GetTasks()
	if ts != nil {
		if task, ok := ts[id]; ok {
			if len(task.ID()) == 0 {
				task.Metadata = &ObjectMetadata{
					Id:        id,
					CreatedAt: m.Metadata.CreatedAt,
				}
			}
			if task.Spec == nil {
				task.Spec = m.GetSpec().TaskSpec(id)
			}
			return task, ok
		}
	}

	// If not available, try to fetch it from the spec
	if spec := m.GetSpec().TaskSpec(id); spec != nil {
		return &Task{
			Metadata: &ObjectMetadata{
				Id:        id,
				CreatedAt: m.Metadata.CreatedAt,
			},
			Spec: spec,
		}, true
	}

	return nil, false
}

func (m *Workflow) Tasks() map[string]*Task {
	tasks := map[string]*Task{}
	for id := range m.GetStatus().GetTasks() {
		task, _ := m.Task(id)
		tasks[id] = task
	}
	for id := range m.GetSpec().GetTasks() {
		if _, ok := tasks[id]; !ok {
			task, _ := m.Task(id)
			tasks[id] = task
		}
	}
	return tasks
}

//
// WorkflowSpec
//

func (m *WorkflowSpec) TaskIds() []string {
	var ids []string
	for k := range m.Tasks {
		ids = append(ids, k)
	}
	return ids
}

func (m *WorkflowSpec) GetPredecessor() string {
	return m.GetDataflow().GetPredecessor()
}

func (m *WorkflowSpec) SetDescription(s string) *WorkflowSpec {
	m.Description = s
	return m
}

func (m *WorkflowSpec) SetOutput(taskID string) *WorkflowSpec {
	m.OutputTask = taskID
	return m
}

func (m *WorkflowSpec) AddTask(id string, task *TaskSpec) *WorkflowSpec {
	if m.Tasks == nil {
		m.Tasks = map[string]*TaskSpec{}
	}
	m.Tasks[id] = task
	return m
}

func (m *WorkflowSpec) TaskSpec(taskID string) *TaskSpec {
	tasks := m.GetTasks()
	if tasks == nil {
		return nil
	}
	return tasks[taskID]
}

//
// WorkflowStatus
//

func (m *WorkflowStatus) Ready() bool {
	return m.GetStatus() == WorkflowStatus_READY
}

func (m *WorkflowStatus) Failed() bool {
	return m.GetStatus() == WorkflowStatus_FAILED
}

func (m *WorkflowStatus) AddTask(id string, t *Task) {
	if m.Tasks == nil {
		m.Tasks = map[string]*Task{}
	}
	m.Tasks[id] = t
}

//Consent Status

func (cs *ConsentStatus) Permitted() bool {
	stat := cs.GetStatus()
	if stat == ConsentStatus_REVOKED || stat == ConsentStatus_PAUSED {
		return false
	}
	return true
}

//
// Zone
//
var zoneRegexp *regexp.Regexp = nil

func GetZoneRegexp() *regexp.Regexp {
	// a quick hack to avoid constantly rebuilding the regexp string
	if zoneRegexp != nil {
		return zoneRegexp
	}

	regexpStr := "("
	i := 0
	for v, zn := range Zone_name {
		zone := ""
		// skip compiling UNDEF zone into the regex
		if v == 0 {
			continue
		}

		// dont add '|'' regex alterative for first value
		if i != 0 {
			zone = zone + "|"
		} else {
			i = 1
		}

		// fission doesnt allow uppercase characters in environemnt
		// names as such we convert to lowercase
		zone += fmt.Sprintf("-%s", strings.ToLower(zn))
		regexpStr += zone
	}
	regexpStr += ")"
	//compile the regexp string
	zoneRegexp = regexp.MustCompile(regexpStr)
	return zoneRegexp
}

// Given a function-id return a list of zone suffixed function id's
func GenZoneVariants(fnId string) []string {
	fns := []string{}
	for i, v := range Zone_name {
		if i != 0 {
			fns = append(fns, fmt.Sprintf("%s-%s", fnId, strings.ToLower(v)))
		}
	}
	return fns
}

func ZoneVariant(fnRef string, z Zone) string {
	zoneName := Zone_name[int32(z)]
	return fmt.Sprintf("%s-%s", fnRef, strings.ToLower(zoneName))
}

func GetZoneString(fnRef string) string {
	re := GetZoneRegexp()
	zone := re.Find([]byte(fnRef))
	if zone == nil {
		return "UNDEF"
	}
	// convert to string and remove the '-' from begining of string
	return strings.ToUpper(string(zone)[1:])
}

func GetZone(fnRef string) Zone {
	return Zone(Zone_value[GetZoneString(fnRef)])
}

func ParseZoneHints(inpt map[string]interface{}) map[string]Zone {
	taskZones := map[string]Zone{}
	for task, zone := range inpt {
		if zName, ok := zone.(string); ok {
			if z := Zone(Zone_value[zName]); z != Zone_UNDEF {
				taskZones[task] = z
			}
		} else {
			logrus.Errorf("Cannot convert interface to string")
		}
	}
	return taskZones
}

//
// FnRef
//
func (ref *FnRef) GenZone() Zone {
	return GetZone(ref.GetID())
}
