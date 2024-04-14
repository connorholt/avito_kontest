package assert

import (
	"testing"
)

func Equal[T comparable](t *testing.T, value, expectedValue T) {
	t.Helper()

	if value != expectedValue {
		t.Errorf("got: %v; want: %v", value, expectedValue)
	}
}
