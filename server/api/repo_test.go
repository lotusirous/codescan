package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/mock"
)

func TestValidateGithubURL(t *testing.T) {

	cases := []struct {
		Name    string
		URL     string
		WantErr bool
	}{
		{
			Name:    "valid-url",
			URL:     "https://github.com/octocat/Hello-World",
			WantErr: false,
		},
		{
			Name:    "no-proto",
			URL:     "github.com/octocat/Hello-World",
			WantErr: true,
		},
		{
			Name:    "malformed-url",
			URL:     "http//github.com/octocat/Hello-World",
			WantErr: true,
		},
		{
			Name:    "missing-repo-name",
			URL:     "http//github.com/octocat",
			WantErr: true,
		},
		{
			Name:    "long-url",
			URL:     "https://github.com/octocat/Hello-World/very/long/path/in/url",
			WantErr: false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			err := validateURL(test.URL)
			if err != nil {
				if !test.WantErr {
					t.Error(err)
				}
			}
		})
	}

}

var mockRepo = &core.Repository{
	ID:      1,
	HttpURL: "https://github.com/octocat/hello-worId",
	Name:    "hello-world",
}

func TestHandleCreateRepo(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	repos := mock.NewMockRepositoryStore(controller)
	repos.EXPECT().Create(gomock.Any(), gomock.Any()).Do(func(_ context.Context, in *core.Repository) error {
		if got, want := in.HttpURL, "https://github.com/octocat/hello-worId"; got != want {
			t.Errorf("Want repo url %s, got %s", want, got)
		}
		if in.Name == "" {
			t.Errorf("Expect repo name not empty")
		}
		return nil
	})

	in := new(bytes.Buffer)
	json.NewEncoder(in).Encode(mockRepo)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", in)

	HandleCreateRepo(repos)(w, r)
	if got, want := w.Code, 200; want != got {
		t.Errorf("Want response code %d, got %d", want, got)
	}

	out := new(core.Repository)
	json.NewDecoder(w.Body).Decode(out)
	if got, want := out.Name, "hello-world"; got != want {
		t.Errorf("Want repo name %s, got %s", want, got)
	}
	if got := out.Created; got == 0 {
		t.Errorf("Want repo created set to current unix timestamp, got %v", got)
	}
	if got := out.Updated; got == 0 {
		t.Errorf("Want repo updated set to current unix timestamp, got %v", got)
	}
}
