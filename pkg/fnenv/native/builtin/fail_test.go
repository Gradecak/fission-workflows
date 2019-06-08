package builtin

import (
	"testing"

	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/gradecak/fission-workflows/pkg/types/typedvalues"
	"github.com/stretchr/testify/assert"
)

func TestFunctionFail_InvokeEmpty(t *testing.T) {
	fn := &FunctionFail{}
	out, err := fn.Invoke(&types.TaskInvocationSpec{})
	assert.Nil(t, out)
	assert.EqualError(t, err, typedvalues.MustUnwrap(defaultErrMsg).(string))

}

func TestFunctionFail_InvokeString(t *testing.T) {
	fn := &FunctionFail{}
	errMsg := "custom error message"
	out, err := fn.Invoke(&types.TaskInvocationSpec{
		Inputs: types.Input(errMsg),
	})
	assert.Nil(t, out)
	assert.EqualError(t, err, errMsg)
}
