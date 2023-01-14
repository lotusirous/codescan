package core

import (
	"context"
	"time"
)

// Status of a scanning job.
const (
	StatusQueued     = "Queued"
	StatusInProgress = "In Progress"
	StatusSuccess    = "Success"
	StatusFailure    = "Failure"
)

// Scan represents the scan for on a repository.
type Scan struct {
	ID         int64      `json:"id"`
	Repository Repository `json:"repository"`
	Status     string     `json:"status"` // refer to status job.
	// unix timestamp
	EnqueuedAt int64 `json:"enqueuedAt"`
	StartedAt  int64 `json:"startedAt"`
	FinishedAt int64 `json:"finishedAt"`
}

// IsDone returns true if the scan has a completed state.
func (s *Scan) IsDone() bool {
	switch s.Status {
	case StatusSuccess, StatusFailure:
		return true
	default:
		return false
	}
}

// IsFailed returns true if the scan has failed
func (s *Scan) IsFailed() bool {
	return s.Status == StatusFailure
}

// ScanStore defines operations for working with scans.
type ScanStore interface {
	// Update stores the status in the datastore.
	Update(ctx context.Context, s *Scan) error

	// Creates persists a scan in the datastore.
	Create(ctx context.Context, s *Scan) error

	// List returns a list of scans from the datastore.
	List(context.Context) ([]*Scan, error)

	// Count returns a count of scans.
	Count(ctx context.Context) (int64, error)

	// Delete deletes a scan from the datastore.
	Delete(context.Context, *Scan) error
}

type (
	// ScanResult represents the scanned result.
	ScanResult struct {
		ID             string    `json:"id"`
		Status         string    `json:"status"`
		RepositoryName string    `json:"repositoryName"`
		RepositoryURL  string    `json:"repositoryURL"`
		Findings       []Finding `json:"findings"`
		EnqueuedAt     time.Time `json:"enqueuedAt"`
		StartedAt      time.Time `json:"startedAt"`
		FinishedAt     time.Time `json:"finishedAt"`
	}
)

type Finding struct {
	Type     string   `json:"type"`
	RuleID   string   `json:"ruleId"`
	Location Location `json:"location"`
	Metadata Metadata `json:"metadata"`
}

type Location struct {
	Path      string    `json:"path"`
	Positions Positions `json:"positions"`
}

type Positions struct {
	Begin Begin `json:"begin"`
}

type Begin struct {
	Line int64 `json:"line"`
}

type Metadata struct {
	Description string `json:"description"`
	Severity    string `json:"severity"`
}
