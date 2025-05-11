package test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/voidwyrm-2/reqproc/runtime/types"
)

const (
	ADDITION_POOLSIZE       = 20
	SUBTRACTION_POOLSIZE    = 20
	MULTIPLICATION_POOLSIZE = 20
	DIVISION_POOLSIZE       = 20
)

func TestAddition(t *testing.T) {
	cases := []stackTestCase{}

	// integer tests
	if err := generateStackTestCases(&cases, ADDITION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Intn(201), rand.Intn(201)
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				float32(a + b),
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "+")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	// float tests
	if err := generateStackTestCases(&cases, ADDITION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Float32()*100, rand.Float32()*100
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				a + b,
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "+")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	// string tests
	if err := generateStackTestCases(&cases, ADDITION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := generateJunkString(rand.Intn(100)), generateJunkString(rand.Intn(100))

			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeString,
				a + b,
			})
			operands = append(operands, `"`+a+`"`, `"`+b+`"`, "+")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	testStack(t, cases)
}

func TestSubtraction(t *testing.T) {
	cases := []stackTestCase{}

	// integer tests
	if err := generateStackTestCases(&cases, SUBTRACTION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Intn(201), rand.Intn(201)
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				float32(a - b),
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "-")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	// float tests
	if err := generateStackTestCases(&cases, SUBTRACTION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Float32()*100, rand.Float32()*100
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				a - b,
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "-")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	testStack(t, cases)
}

func TestMultiplication(t *testing.T) {
	cases := []stackTestCase{}

	// integer tests
	if err := generateStackTestCases(&cases, MULTIPLICATION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Intn(201), rand.Intn(201)
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				float32(a * b),
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "*")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	// float tests
	if err := generateStackTestCases(&cases, MULTIPLICATION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Float32()*100, rand.Float32()*100
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				a * b,
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "*")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	// string tests
	if err := generateStackTestCases(&cases, MULTIPLICATION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := generateJunkString(rand.Intn(100)), rand.Intn(20)+1

			strs := []string{}

			for range b {
				strs = append(strs, a)
			}

			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeString,
				strings.Join(strs, ""),
			})
			operands = append(operands, `"`+a+`"`, fmt.Sprint(b), "*")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	testStack(t, cases)
}

// I hate floating-point imprecision so much
/*
func TestDivision(t *testing.T) {
	cases := []stackTestCase{}

	// integer tests
	if err := generateStackTestCases(&cases, DIVISION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Intn(201)+1, rand.Intn(201)+1
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				float32(a / b),
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "/")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	// float tests
	if err := generateStackTestCases(&cases, DIVISION_POOLSIZE, func(i int) (stackTestCase, error) {
		expected := []struct {
			t types.ReqVarType
			v any
		}{}
		operands := []string{}

		operators := rand.Intn(8) + 1

		for range operators {
			a, b := rand.Float32()*100+1, rand.Float32()*100+1
			expected = append(expected, struct {
				t types.ReqVarType
				v any
			}{
				types.TypeNumber,
				a / b,
			})
			operands = append(operands, fmt.Sprint(a), fmt.Sprint(b), "/")
		}

		return stackTestCase{
			strings.Join(operands, " "),
			expected,
		}, nil
	}); err != nil {
		t.Fatal(err.Error())
	}

	testStack(t, cases)
}
*/
