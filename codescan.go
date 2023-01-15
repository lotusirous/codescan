package codescan

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lotusirous/codescan/checker"
	"github.com/lotusirous/codescan/config"
	"github.com/lotusirous/codescan/core"
	"github.com/lotusirous/codescan/fetcher/github"
	"github.com/lotusirous/codescan/sched"
	"github.com/lotusirous/codescan/server"
	"github.com/lotusirous/codescan/signal"
	"github.com/lotusirous/codescan/store/db"
	"github.com/lotusirous/codescan/store/repos"
	"github.com/lotusirous/codescan/store/scans"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Run starts codescan program.
func Run() error {
	conf, err := config.Environ()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	setupLogger(conf.Debug)

	db, err := db.Connect(conf.Database.Datasource, conf.Database.MaxConnections)
	if err != nil {
		return err
	}

	ctx := signal.WithContext(context.Background())
	var (
		scanStore  = scans.New(db)
		repoStore  = repos.New(db)
		gitFetcher = github.New("tmp", "codescan")

		manager = sched.New(4, scanStore,
			repoStore,
			gitFetcher,
			&core.Scanner{Type: "sast", Analyzers: checker.DefaultRules()},
		)

		srv = server.New(
			conf.ServerAddress,
			repoStore,
			scanStore,
			manager,
		)
	)

	// Handle restore last scan from db.
	if err := manager.RestoreLastScan(ctx); err != nil {
		return err
	}

	var g errgroup.Group
	g.Go(func() error {
		log.Info().Str("addr", conf.ServerAddress).Msg("server started")
		return srv.ListenAndServe(ctx)
	})

	g.Go(func() error {
		return manager.Start(ctx)
	})

	return g.Wait()

}

func setupLogger(debug bool) {
	w := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	if debug {
		log.Logger = log.Output(w).Level(zerolog.DebugLevel)
		return
	}
	log.Logger = log.Output(w).Level(zerolog.InfoLevel)
}
