package errors

import "testing"

func TestError(t *testing.T) {
	got, want := ErrNotFound.Error(), ErrNotFound.(*Error).Message
	if got != want {
		t.Errorf("Want error string %q, got %q", got, want)
	}
}
