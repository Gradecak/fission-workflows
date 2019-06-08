package builtin

import (
	"github.com/gradecak/fission-workflows/pkg/types"
	"github.com/gradecak/fission-workflows/pkg/types/typedvalues"
	"github.com/sirupsen/logrus"
)

const (
	Compose      = "compose"
	ComposeInput = types.InputMain
)

/*
FunctionCompose provides a way to merge, modify and create complex values from multiple inputs.
Other than outputting the composed inputs, compose does not perform any other operation.
This is useful when you want to merge the outputs from different tasks (for example in a MapReduce or scatter-gather
scenario).

**Specification**

**input**   | required | types  | description
------------|----------|--------|---------------------------------
default     | no       | *      | the inputs to be merged into a single map or outputted if none other.
*           | no       | *      | the inputs to be merged into a single map.

**Note: custom message does not yet propagate back to the user**

**output** (*) The composed map, single default input, or nothing.

**Example**

Compose with a single input, similar to `noop`:
```yaml
# ...
foo:
  run: compose
  inputs: "all has failed"
# ...
```

Composing a map inputs:
```yaml
# ...
foo:
  run: compose
  inputs:
    foo: bar
    fission: workflows
# ...
```
*/
// TODO avoid adding function-injected fields to compose
type FunctionCompose struct{}

func (fn *FunctionCompose) Invoke(spec *types.TaskInvocationSpec) (*typedvalues.TypedValue, error) {

	var output *typedvalues.TypedValue
	switch len(spec.GetInputs()) {
	case 0:
		output = nil
	case 1:
		defaultInput, ok := spec.GetInputs()[ComposeInput]
		if ok {
			output = defaultInput
			break
		}
		fallthrough
	default:
		results := map[string]interface{}{}
		for k, v := range spec.GetInputs() {
			i, err := typedvalues.Unwrap(v)
			if err != nil {
				return nil, err
			}
			results[k] = i
		}
		p, err := typedvalues.Wrap(results)
		if err != nil {
			return nil, err
		}
		output = p
	}
	logrus.Infof("[internal://%s] %v (Type: %s, Labels: %v)", Compose, typedvalues.MustUnwrap(output), output.ValueType(),
		output.GetMetadata())
	return output, nil
}
