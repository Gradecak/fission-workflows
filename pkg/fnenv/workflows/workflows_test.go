package workflows

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/gradecak/fission-workflows/pkg/api"
	"github.com/gradecak/fission-workflows/pkg/api/projectors"
	"github.com/gradecak/fission-workflows/pkg/api/store"
	"github.com/gradecak/fission-workflows/pkg/fes"
	"github.com/gradecak/fission-workflows/pkg/fes/backend/mem"
	"github.com/gradecak/fission-workflows/pkg/fes/cache"
	"github.com/gradecak/fission-workflows/pkg/fes/testutil"
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/gradecak/fission-workflows/pkg/types/typedvalues"
	"github.com/gradecak/fission-workflows/pkg/util"
	"github.com/stretchr/testify/assert"
)

const (
	workflowID     = "123"
	defaultTimeout = 5 * time.Second
)

func defaultDeadline() time.Time {
	return time.Now().Add(defaultTimeout)
}

func TestRuntime_InvokeWorkflow_SubTimeout(t *testing.T) {
	runtime, _, _, _ := setup()
	spec := types.NewWorkflowInvocationSpec(workflowID, time.Now().Add(10*time.Millisecond))
	_, err := runtime.InvokeWorkflow(spec)
	assert.Equal(t, api.ErrInvocationCanceled, err.Error())
}

func TestRuntime_InvokeWorkflow_PollTimeout(t *testing.T) {
	runtime, _, _, _ := setup()
	runtime.invocations = store.NewInvocationStore(testutil.NewCache()) // ensure that cache does not support pubsub
	runtime.pollInterval = 10 * time.Millisecond
	_, err := runtime.InvokeWorkflow(types.NewWorkflowInvocationSpec(workflowID, defaultDeadline()))
	assert.EqualError(t, err, context.DeadlineExceeded.Error())
}

func TestRuntime_InvokeWorkflow_InvalidSpec(t *testing.T) {
	runtime, _, _, _ := setup()
	_, err := runtime.InvokeWorkflow(types.NewWorkflowInvocationSpec("", defaultDeadline()))
	assert.IsType(t, errors.New(""), err)
}

func TestRuntime_InvokeWorkflow_SubSuccess(t *testing.T) {
	runtime, invocationAPI, _, cc := setup()
	output := typedvalues.MustWrap("foo")
	outputHeaders := typedvalues.MustWrap(typedvalues.MustWrap(map[string]interface{}{
		"some-key": "some-value",
	}))
	go func() {
		// Simulate workflow invocation
		time.Sleep(50 * time.Millisecond)
		entities := cc.List()
		wfiID := entities[0].Id
		err := invocationAPI.Complete(wfiID, output, outputHeaders)
		if err != nil {
			panic(err)
		}
	}()
	wfi, err := runtime.InvokeWorkflow(types.NewWorkflowInvocationSpec(workflowID, defaultDeadline()))
	assert.NoError(t, err)
	util.AssertProtoEqual(t, output, wfi.GetStatus().GetOutput())
	util.AssertProtoEqual(t, outputHeaders, wfi.GetStatus().GetOutputHeaders())
	assert.True(t, wfi.GetStatus().Finished())
	assert.True(t, wfi.GetStatus().Successful())
}

func TestRuntime_InvokeWorkflow_PollSuccess(t *testing.T) {
	runtime, invocationAPI, _, c := setup()
	pollCache := store.NewInvocationStore(testutil.NewCache()) // ensure that cache does not support pubsub
	runtime.invocations = pollCache

	output := typedvalues.MustWrap("foo")
	outputHeaders := typedvalues.MustWrap(typedvalues.MustWrap(map[string]interface{}{
		"some-key": "some-value",
	}))
	go func() {
		// Simulate workflow invocation
		time.Sleep(50 * time.Millisecond)
		entities := c.List()
		wfiID := entities[0].Id
		err := invocationAPI.Complete(wfiID, output, outputHeaders)
		if err != nil {
			panic(err)
		}
		time.Sleep(50 * time.Millisecond)
		entity, err := c.GetAggregate(projectors.NewInvocationAggregate(wfiID))
		assert.NoError(t, err)
		pollCache.CacheReader.(fes.CacheReaderWriter).Put(entity)
	}()
	wfi, err := runtime.InvokeWorkflow(types.NewWorkflowInvocationSpec(workflowID, defaultDeadline()))
	assert.NoError(t, err)
	util.AssertProtoEqual(t, output, wfi.GetStatus().GetOutput())
	util.AssertProtoEqual(t, outputHeaders, wfi.GetStatus().GetOutputHeaders())
	assert.True(t, wfi.GetStatus().Finished())
	assert.True(t, wfi.GetStatus().Successful())
}

func TestRuntime_InvokeWorkflow_Fail(t *testing.T) {
	runtime, invocationAPI, _, cache := setup()
	wfiErr := errors.New("stub err")
	go func() {
		// Simulate workflow invocation
		time.Sleep(50 * time.Millisecond)
		entities := cache.List()
		wfiID := entities[0].Id
		err := invocationAPI.Fail(wfiID, wfiErr)
		if err != nil {
			panic(err)
		}
	}()
	wfi, err := runtime.InvokeWorkflow(types.NewWorkflowInvocationSpec(workflowID, defaultDeadline()))
	assert.NoError(t, err)
	assert.Equal(t, wfiErr.Error(), wfi.GetStatus().GetError().Error())
	assert.True(t, wfi.GetStatus().Finished())
	assert.False(t, wfi.GetStatus().Successful())
}

func TestRuntime_InvokeWorkflow_Cancel(t *testing.T) {
	runtime, invocationAPI, _, cache := setup()
	go func() {
		// Simulate workflow invocation
		time.Sleep(50 * time.Millisecond)
		entities := cache.List()
		wfiID := entities[0].Id
		err := invocationAPI.Cancel(wfiID)
		if err != nil {
			panic(err)
		}
	}()
	wfi, err := runtime.InvokeWorkflow(types.NewWorkflowInvocationSpec(workflowID, defaultDeadline()))
	assert.NoError(t, err)
	assert.Equal(t, api.ErrInvocationCanceled, wfi.GetStatus().GetError().Error())
	assert.True(t, wfi.GetStatus().Finished())
	assert.False(t, wfi.GetStatus().Successful())
}

func TestRuntime_Invoke(t *testing.T) {
	runtime, invocationAPI, _, cache := setup()

	deadline, _ := ptypes.TimestampProto(time.Now().Add(10 * time.Second))
	fnref := types.NewFnRef("workflows", "", workflowID)
	spec := types.NewTaskInvocationSpec(&types.WorkflowInvocation{
		Metadata: types.NewObjectMetadata("wi-123"),
		Spec: &types.WorkflowInvocationSpec{
			Deadline: deadline,
		},
	}, &types.Task{
		Metadata: types.NewObjectMetadata("ti-123"),
		Spec:     &types.TaskSpec{},
		Status: &types.TaskStatus{
			FnRef: &fnref,
		},
	}, time.Now())
	spec.Inputs = types.Inputs{
		types.InputParent: typedvalues.MustWrap("parentID"),
	}
	output := typedvalues.MustWrap("foo")
	outputHeaders := typedvalues.MustWrap(typedvalues.MustWrap(map[string]interface{}{
		"some-key": "some-value",
	}))
	go func() {
		// Simulate workflow invocation
		time.Sleep(50 * time.Millisecond)
		entities := cache.List()
		wfiID := entities[0].Id
		err := invocationAPI.Complete(wfiID, output, outputHeaders)
		if err != nil {
			panic(err)
		}
	}()

	task, err := runtime.Invoke(spec)
	assert.NoError(t, err)
	util.AssertProtoEqual(t, output, task.GetOutput())
	util.AssertProtoEqual(t, outputHeaders, task.GetOutputHeaders())
}

func setup() (*Runtime, *api.Invocation, *mem.Backend, fes.CacheReaderWriter) {
	backend := mem.NewBackend()
	invocationAPI := api.NewInvocationAPI(backend)
	workflowsCache := testutil.NewCache()
	err := workflowsCache.Put(&types.Workflow{
		Metadata: &types.ObjectMetadata{
			Id: workflowID,
		},
		Status: &types.WorkflowStatus{
			Status: types.WorkflowStatus_READY,
		},
	})
	if err != nil {
		panic(err)
	}
	workflows := store.NewWorkflowsStore(workflowsCache)
	c := cache.NewSubscribedCache(testutil.NewCache(), projectors.NewWorkflowInvocation(), backend.Subscribe())
	runtime := NewRuntime(invocationAPI, store.NewInvocationStore(c), workflows)
	return runtime, invocationAPI, backend, c
}
