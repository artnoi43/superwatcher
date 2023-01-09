package testutils

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

var (
	mapTypeFuncFlagVar = map[string]interface{}{
		"string": flag.BoolVar,
		"bool":   flag.BoolVar,
		"int":    flag.IntVar,
	}

	flagCase    int
	flagVerbose bool
)

func init() {
	RegFlagVar(&flagCase, "case", -1, "test case")
	RegFlagVar(&flagVerbose, "verbose", false, "test output verbosity")

	flagCase = GetFlagVar[int]("case")
	flagVerbose = GetFlagVar[bool]("verbose")
}

func RegFlagVar[T any](p *T, name string, value T, usage string) {
	if flag.Lookup(name) == nil {
		t := reflect.TypeOf(value).String()
		f, ok := mapTypeFuncFlagVar[t]
		if !ok {
			panic("RegFlagVar does not support type " + t)
		}

		f.(func(*T, string, T, string))(p, name, value, usage)
	}
}

func GetFlagVar[T any](name string) T {
	return flag.Lookup(name).Value.(flag.Getter).Get().(T)
}

// TODO: verbose does not work
func GetFlagValues() (caseNumber int, verbose bool) {
	caseNumber = 1 // default case 1
	if flagCase != 0 {
		caseNumber = flagCase
	}
	if flagVerbose {
		verbose = flagVerbose
	}

	return caseNumber, verbose
}

func CheckTestCase[T any](t *testing.T, caseNumber int, cases []T) {
	if l := len(cases); caseNumber >= l {
		t.Skipf("no such test case %d", caseNumber)
	}
}

func RunTestCase[T any](
	t *testing.T,
	testName string,
	testCases []T,
	testFunc func(*testing.T, int) error,
) error {
	caseNumber, _ := GetFlagValues()
	CheckTestCase(t, caseNumber, testCases)

	testName += " case %d"

	var err error
	if caseNumber > 0 {
		testName = fmt.Sprintf(testName, caseNumber)
		t.Run(testName, func(t *testing.T) {
			err = testFunc(t, caseNumber)
		})

		if err != nil {
			return errors.Wrapf(err, "%s returned an error", testName)
		}
	}

	for i := range testCases {
		caseNumber = i + 1
		testName = fmt.Sprintf(testName, caseNumber)
		t.Run(testName, func(t *testing.T) {
			err = testFunc(t, caseNumber)
		})

		if err != nil {
			return errors.Wrapf(err, "%s returned an error", testName)
		}
	}

	return nil
}
