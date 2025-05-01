package test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/voidwyrm-2/reqproc/runtime/interpreter"
	"github.com/voidwyrm-2/reqproc/runtime/scope"
	"github.com/voidwyrm-2/reqproc/runtime/types"
)

type testCase[I, E any] struct {
	input    I
	expected []E
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
	for _, c := range cases {
		i, err := interpreter.New(scope.New(nil, map[string]types.ReqType{}))
		if err != nil {
			t.Fatal(err.Error())
		}

		result, err := i.Execute(c.input)
		if err != nil {
			t.Fatal(err.Error())
		}

		if len(result) != len(c.expected) {
			fmt.Println(result, c.input)
			t.Fatalf("expected %d values to be on the stack, but found %d instead", len(c.expected), len(result))
		}

		for i, j := len(result)-1, 0; j < len(c.expected); j++ {
			r := result[i]
			e := c.expected[j]

			// the test reults were consistently getting a .00001 difference
			cond := r.Literal() == e.v
			if f1, ok := r.Literal().(float32); ok {
				if f2, ok := e.v.(float32); ok {
					cond = f1 >= f2-.00001 && f1 <= f2+.00001
				}
			}

			if r.Type() != e.t {
				t.Fatalf("expected type %s at position %d, but found %s instead", e.t.String(), i, r.Type().String())
			} else if !cond {
				t.Fatalf("expected value '%v' at position %d, but found '%v' instead", e.v, j, r.Literal())
			}

			i--
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
