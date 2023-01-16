package sched

import (
	"context"
	"time"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/core"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

type msg struct {
	scan *core.Scan
	repo *core.Repository
}

func New(
	workers int,
	scans core.ScanStore,
	repos core.RepositoryStore,
	results core.ScanResultStore,
	git core.GitFetcher,
	scanner core.Scanner,
) core.ScanScheduler {
	return &scheduler{
		workers:     workers,
		Scans:       scans,
		Git:         git,
		Repos:       repos,
		ScanResults: results,
		Scanner:     scanner,
		Queue:       make(chan msg, 1000),
	}
}

type scheduler struct {
	workers     int
	Scans       core.ScanStore
	Git         core.GitFetcher
	Repos       core.RepositoryStore
	ScanResults core.ScanResultStore
	Scanner     core.Scanner
	Queue       chan msg
}

// RestoreLastJob get the enqueue job from datastore and put it the queue.
func (s *scheduler) RestoreLastScan(ctx context.Context) error {
	scans, err := s.Scans.FindEnqueued(ctx)
	if err != nil {
		return err
	}
	for _, scan := range scans {
		repo, err := s.Repos.Find(ctx, scan.RepoID)
		if err != nil {
			return err
		}
		s.Queue <- msg{repo: repo, scan: scan}
	}
	return nil
}

func (s *scheduler) ScanRepo(ctx context.Context, repo *core.Repository) (*core.Scan, error) {
	scan := &core.Scan{
		RepoID:     repo.ID,
		Status:     core.StatusQueued,
		EnqueuedAt: time.Now().Unix(),
	}
	if err := s.Scans.Create(ctx, scan); err != nil {
		return nil, err
	}
	s.Queue <- msg{repo: repo, scan: scan}
	return scan, nil
}

func (s *scheduler) Start(ctx context.Context) error {
	var g errgroup.Group
	for i := 0; i < s.workers; i++ {
		g.Go(func() error {
			return s.doWork(ctx)
		})
	}
	log.Info().Msgf("scan by %d workers", s.workers)
	return g.Wait()
}

func (s *scheduler) doWork(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m := <-s.Queue:
			err := s.work(ctx, m.scan, m.repo)
			if err != nil {
				log.Warn().Err(err).Msg("do work failed")
			}
		}
	}
}

func (s *scheduler) workWithStatus(ctx context.Context, job *core.Scan, fn func() error) error {

	job.Status = core.StatusInProgress
	job.StartedAt = time.Now().Unix()
	if err := s.Scans.Update(ctx, job); err != nil {
		return err
	}

	if err := fn(); err != nil {
		job.Status = core.StatusFailure
		log.Error().Err(err).Msg("work fn failed")
		return s.Scans.Update(ctx, job)
	}
	job.Status = core.StatusSuccess
	return s.Scans.Update(ctx, job)
}

func (s *scheduler) work(ctx context.Context, job *core.Scan, repo *core.Repository) error {

	return s.workWithStatus(ctx, job, func() error {
		dir, cleanup, err := s.Git.Clone(repo.HttpURL)
		if err != nil {
			return err
		}
		defer func() {
			if err := cleanup(); err != nil {
				log.Warn().Err(err).Msg("clean up directory failed")
			}
		}()
		summary, err := s.Git.Summarize(dir)
		if err != nil {
			return err
		}
		log.Info().Msg("done summary repo")

		diags, err := s.Scanner.Scan(dir)
		if err != nil {
			return err
		}
		findings := s.toFindings(diags)
		return s.ScanResults.Create(ctx, &core.ScanResult{
			ScanID:   job.ID,
			RepoID:   repo.ID,
			Commit:   summary.CommitHash,
			Created:  time.Now().Unix(),
			Updated:  time.Now().Unix(),
			Findings: findings,
		})

	})

}

func (s *scheduler) toFindings(diags []*analysis.Diagnostic) []core.Finding {
	out := make([]core.Finding, 0)
	for _, diag := range diags {
		out = append(out, core.Finding{
			Type:   s.Scanner.Type,
			RuleID: diag.ByAnalyzer.Name,
			Location: core.Location{
				Path:      diag.Path,
				Positions: core.Positions{Begin: core.Begin{Line: int64(diag.Pos)}},
			},
			Metadata: core.Metadata{
				Description: diag.ByAnalyzer.Meta.Description,
				Severity:    diag.ByAnalyzer.Meta.Severity,
			},
		})
	}
	return out
}
