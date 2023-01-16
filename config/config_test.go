package config

import "testing"

func TestEnviron(t *testing.T) {
	conf, err := Environ()
	if err != nil {
		t.Error(err)
	}

	if want, got := ":8080", conf.ServerAddress; got != want {
		t.Errorf("Want default address %s - got %s", want, got)
	}
}
