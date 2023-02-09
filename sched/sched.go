package sched

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/lotusirous/codescan/checker/analysis"
	"github.com/lotusirous/codescan/core"
	"github.com/rs/zerolog/log"
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
		log.Debug().Int64("repo_id", repo.ID).Int64("scan_id", scan.ID).Msg("restore last scan")
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

// Loop starts the worker in the group
// workers should not stop even if the error is occurred.
func (s *scheduler) Loop(ctx context.Context) {
	var wg sync.WaitGroup
	for i := 0; i < s.workers; i++ {
		wg.Add(1)
		func() {
			if err := s.workContext(ctx); err != nil {
				log.Warn().Err(err).Msg("workContext failed")
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func (s *scheduler) workContext(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m := <-s.Queue:
			err := s.analyze(ctx, m.scan, m.repo)
			if errors.Is(err, context.Canceled) {
				return nil
			}
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
	job.FinishedAt = time.Now().Unix()
	return s.Scans.Update(ctx, job)
}

func (s *scheduler) analyze(ctx context.Context, job *core.Scan, repo *core.Repository) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
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
			log.Info().Int64("scan_id", job.ID).Str("repo", repo.HttpURL).Msg("start scanning repo")
			defer func() {
				log.Info().Int64("scan_id", job.ID).Str("repo", repo.HttpURL).Msg("finished scanning repo")
			}()

			diags, err := s.Scanner.Scan(dir)
			if err != nil {
				return err
			}
			findings := s.toFindings(diags, dir)

			// the store should not use the parent context
			// because it might be cancelled by the signal.
			return s.ScanResults.Create(context.Background(), &core.ScanResult{
				ScanID:   job.ID,
				RepoID:   repo.ID,
				Commit:   summary.CommitHash,
				Created:  time.Now().Unix(),
				Updated:  time.Now().Unix(),
				Findings: findings,
			})
		})
	}
}

func (s *scheduler) toFindings(diags []*analysis.Diagnostic, stripDir string) []core.Finding {
	out := make([]core.Finding, 0)
	for _, diag := range diags {
		out = append(out, core.Finding{
			Type:   s.Scanner.Type,
			RuleID: diag.ByAnalyzer.Name,
			Location: core.Location{
				Path:      strings.ReplaceAll(diag.Path, stripDir, ""),
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
