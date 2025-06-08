package test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/voidwyrm-2/reqproc/runtime/interpreter"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type testCase[I, E any] struct {
	input    I
	expected []E
	log      bool
}

type stackTestCase = testCase[string, struct {
	t types.ReqVarType
	v any
}]

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// sourced from: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go#31832326
func generateJunkString(length int) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}

	return string(b)
}

func testStack(t *testing.T, cases []stackTestCase) {
	for caseIndex, c := range cases {
		if c.log {
			t.Logf("(%d of %d) testing `%s`\n", caseIndex+1, len(cases), c.input)
		} else {
			t.Logf("testing %d of %d cases\n", caseIndex+1, len(cases))
		}

		i, err := interpreter.New(scope.New(nil, map[string]types.ReqType{}))
		if err != nil {
			t.Error(err.Error())
		}

		result, err := i.Execute(c.input)
		if err != nil {
			t.Error(err.Error() + " (with `" + c.input + "`)")
		}

		if len(result) != len(c.expected) {
			t.Errorf("expected %d values to be on the stack, but found %d instead (with `%s`)", len(c.expected), len(result), c.input)
		}

		for i := len(result) - 1; i != 0; i-- {
			e, r := c.expected[i], result[i]

			cond := reflect.DeepEqual(r.Literal(), e.v)
			if e.v == nil {
				cond = true
			} else if f1, ok := r.Literal().(float32); ok { // the test reults were consistently getting a .00001 difference
				if f2, ok := e.v.(float32); ok {
					cond = f1 >= f2-.00001 && f1 <= f2+.00001
				}
			} else if checkf, ok := e.v.(func(types.ReqType) bool); ok {
				cond = checkf(r)
			}

			if r.Type()&e.t != r.Type() {
				t.Errorf("expected type %s at position %d, but found %s instead (with `%s`)", e.t.String(), i, r.Type().String(), c.input)
			} else if !cond {
				t.Errorf("expected value '%v' at position %d, but found '%v' instead (with `%s`)", e.v, i, r.Literal(), c.input)
			}
		}

		if c.log {
			t.Logf("(%d of %d) test output: `%s`\n", caseIndex+1, len(cases), i.GetStack())
		} else {
			t.Logf("tested %d of %d cases\n", caseIndex+1, len(cases))
		}
	}
}

func generateStackTestCases(cases *[]stackTestCase, amount int, generator func(i int) (stackTestCase, error)) error {
	for i := range amount {
		if generated, err := generator(i); err != nil {
			return err
		} else {
			*cases = append(*cases, generated)
		}
	}

	return nil
}
