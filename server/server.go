package server

import (
	"context"
	"html/template"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/server/api"
	"golang.org/x/sync/errgroup"
)

// TemplateSet is a set of template for rendering.
type TemplateSet = map[string]*template.Template

func New(
	addr string,
	repos core.RepositoryStore,
	scans core.ScanStore,
	scheduler core.ScanScheduler,
) Server {
	return Server{
		Addr:  addr,
		Repos: repos,
		Scans: scans,
		Sched: scheduler,
	}
}

// A Server defines parameters for running an HTTP server.
// The TLS will be applied in this struct also.
type Server struct {
	Addr  string
	Sched core.ScanScheduler
	Repos core.RepositoryStore
	Scans core.ScanStore
}

func (s Server) handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Recoverer)
	// r.Get("/", s.handleHome())
	r.Route("/api", func(r chi.Router) {
		r.Get("/repos", api.HandleListRepo(s.Repos))
		r.Post("/repos", api.HandleCreateRepo(s.Repos))
		r.Delete("/repos/{id}", api.HandleDeleteRepo(s.Repos))

		r.Post("/scans", api.HandleListScan(s.Scans))
		r.Post("/scans", api.HandleScanRepo(s.Sched, s.Repos, s.Scans))
		r.Post("/scans/{id}", api.HandleFindScan(s.Scans))
	})
	return r
}

const timeoutGracefulShutdown = 5 * time.Second

// ListenAndServe initializes a server to respond to HTTP network requests.
func (s Server) ListenAndServe(ctx context.Context) error {
	err := s.listenAndServe(ctx)
	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func (s Server) listenAndServe(ctx context.Context) error {
	var g errgroup.Group
	s1 := &http.Server{
		Addr:        s.Addr,
		Handler:     s.handler(),
		ReadTimeout: 60 * time.Second, // magic number from nginx.
	}
	g.Go(func() error {
		<-ctx.Done()

		ctxShutdown, cancel := context.WithTimeout(context.Background(), timeoutGracefulShutdown)
		defer cancel()

		return s1.Shutdown(ctxShutdown)
	})
	g.Go(s1.ListenAndServe)
	return g.Wait()
}

func (s Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Please go to /api"))
	}
}
