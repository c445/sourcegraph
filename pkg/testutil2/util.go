package testutil2

import (
	"reflect"
	"testing"
)

func CheckDeepEqual(t *testing.T, exp, actual interface{}, errMsg string) {
	if !reflect.DeepEqual(exp, actual) {
		t.Errorf("%s: (exp != actual) %+v != %+v", errMsg, exp, actual)
	}
}
