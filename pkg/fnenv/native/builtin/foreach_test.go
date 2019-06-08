package builtin

import (
	"testing"

	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/gradecak/fission-workflows/pkg/types/typedvalues"
	"github.com/gradecak/fission-workflows/pkg/types/typedvalues/controlflow"
	"github.com/stretchr/testify/assert"
)

func TestFunctionForeach_Invoke(t *testing.T) {
	foreachElements := []interface{}{1, 2, 3, 4, "foo"}
	out, err := (&FunctionForeach{}).Invoke(&types.TaskInvocationSpec{
		Inputs: map[string]*typedvalues.TypedValue{
			ForeachInputForeach: typedvalues.MustWrap(foreachElements),
			ForeachInputDo: typedvalues.MustWrap(&types.TaskSpec{
				FunctionRef: Noop,
			}),
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, controlflow.TypeWorkflow, out.ValueType())

	wf, err := controlflow.UnwrapWorkflow(out)
	assert.NoError(t, err)
	assert.Equal(t, len(foreachElements)+1, len(wf.Tasks)) // + 1 for the noop-task in the foreach loop.
	assert.NotNil(t, wf.Tasks["do_0"])
	assert.Equal(t, foreachElements[0], int(typedvalues.MustUnwrap(wf.Tasks["do_0"].Inputs["_item"]).(int32)))
}
