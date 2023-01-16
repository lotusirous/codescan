package render

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/lotusirous/codescan/server/api/errors"
)

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("pc load letter")
	InternalError(w, err)

	if got, want := w.Code, 500; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &errors.Error{}
	if err := json.NewDecoder(w.Body).Decode(errjson); err != nil {
		t.Error("Unable to decode errjson")
	}
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestWriteNotFound(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("pc load letter")
	NotFound(w, err)

	if got, want := w.Code, 404; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &errors.Error{}
	if err := json.NewDecoder(w.Body).Decode(errjson); err != nil {
		t.Error("Unable to decode errjson")
	}
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestWriteNotFoundf(t *testing.T) {
	w := httptest.NewRecorder()

	NotFoundf(w, "pc %s", "load letter")
	if got, want := w.Code, 404; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &errors.Error{}
	if err := json.NewDecoder(w.Body).Decode(errjson); err != nil {
		t.Error("Unable to decode errjson")
	}
	if got, want := errjson.Message, "pc load letter"; got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}

func TestWriteInternalError(t *testing.T) {
	w := httptest.NewRecorder()

	err := errors.New("pc load letter")
	InternalError(w, err)

	if got, want := w.Code, 500; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	errjson := &errors.Error{}
	if err := json.NewDecoder(w.Body).Decode(errjson); err != nil {
		t.Error("Unable to decode errjson")
	}
	if got, want := errjson.Message, err.Error(); got != want {
		t.Errorf("Want error message %s, got %s", want, got)
	}
}
