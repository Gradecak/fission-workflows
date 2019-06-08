package expr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gradecak/fission-workflows/pkg/types/typedvalues"
	"github.com/stretchr/testify/assert"
)

var scope = map[string]interface{}{
	"bit": "bat",
}
var rootScope = map[string]interface{}{
	"foo":          "bar",
	"currentScope": scope,
}

func TestResolveTestRootScopePath(t *testing.T) {

	exprParser := NewJavascriptExpressionParser()

	resolved, err := exprParser.Resolve(rootScope, "", mustParseExpr("{$.currentScope.bit}"))
	if err != nil {
		t.Error(err)
	}

	resolvedString, err := typedvalues.Unwrap(resolved)
	if err != nil {
		t.Error(err)
	}

	expected := scope["bit"]
	assert.Equal(t, expected, resolvedString)
}

func TestResolveTestScopePath(t *testing.T) {
	currentTask := "fooTask"
	exprParser := NewJavascriptExpressionParser()

	resolved, err := exprParser.Resolve(rootScope, currentTask, mustParseExpr("{"+varCurrentTask+"}"))
	assert.NoError(t, err)

	resolvedString, err := typedvalues.Unwrap(resolved)
	assert.NoError(t, err)

	assert.Equal(t, currentTask, resolvedString)
}

func TestResolveLiteral(t *testing.T) {

	exprParser := NewJavascriptExpressionParser()

	expected := "foobar"
	resolved, err := exprParser.Resolve(rootScope, "output", mustParseExpr(fmt.Sprintf("{'%s'}", expected)))
	assert.NoError(t, err)

	resolvedString, _ := typedvalues.Unwrap(resolved)
	assert.Equal(t, expected, resolvedString)
}

func TestResolveTransformation(t *testing.T) {

	exprParser := NewJavascriptExpressionParser()

	src := "foobar"
	expected := strings.ToUpper(src)
	resolved, err := exprParser.Resolve(rootScope, "", mustParseExpr(fmt.Sprintf("{'%s'.toUpperCase()}", src)))
	assert.NoError(t, err)

	resolvedString, _ := typedvalues.Unwrap(resolved)
	assert.Equal(t, expected, resolvedString)
}

func TestResolveInjectedFunction(t *testing.T) {

	exprParser := NewJavascriptExpressionParser()

	resolved, err := exprParser.Resolve(rootScope, "", mustParseExpr("{uid()}"))
	assert.NoError(t, err)

	resolvedString, _ := typedvalues.Unwrap(resolved)

	assert.NotEmpty(t, resolvedString)
}

func mustParseExpr(s string) *typedvalues.TypedValue {
	tv := typedvalues.MustWrap(s)
	if tv.ValueType() != typedvalues.TypeExpression {
		panic(fmt.Sprintf("Should be %v, but was '%v'", typedvalues.TypeExpression, tv.ValueType()))
	}

	return tv
}
