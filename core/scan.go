package core

import (
	"context"

	"github.com/lotusirous/codescan/checker/analysis"
)

// Status of a scanning job.
const (
	StatusQueued     = "Queued"
	StatusInProgress = "In Progress"
	StatusSuccess    = "Success"
	StatusFailure    = "Failure"
)

// Scanner represents a scanner type in the system.
type Scanner struct {
	Type      string // sast (static) or dast (dynamic)
	Analyzers []*analysis.Analyzer
	Scan      func(dir string) ([]*analysis.Diagnostic, error)
}

// Scan represents the scan for on a repository.
type Scan struct {
	ID         int64  `json:"id"`
	RepoID     int64  `json:"repository"`
	Status     string `json:"status"` // refer to status scanning job
	EnqueuedAt int64  `json:"enqueuedAt"`
	StartedAt  int64  `json:"startedAt"`
	FinishedAt int64  `json:"finishedAt"`
	Error      error  // only set when the scan is failed.
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

// ScanStore defines operations for working with scans.
type ScanStore interface {

	// FindEnqueued supports the restored job.
	FindEnqueued(ctx context.Context) ([]*Scan, error)

	// Update stores the status in the datastore.
	Update(ctx context.Context, s *Scan) error

	// Find returns a scan from datastore..
	Find(ctx context.Context, id int64) (*Scan, error)

	// Creates persists a scan in the datastore.
	Create(ctx context.Context, s *Scan) error

	// List returns a list of scans from the datastore.
	List(context.Context) ([]*Scan, error)

	// Count returns a count of scans.
	Count(ctx context.Context) (int64, error)

	// Delete deletes a scan from the datastore.
	Delete(context.Context, *Scan) error
}

// ScanResults is the result of scanning.
type ScanResult struct {
	ID       int64     `json:"id"`
	ScanID   int64     `json:"scan_id"`
	RepoID   int64     `json:"repo_id"`
	Commit   string    `json:"commit"` // latest commit
	Created  int64     `json:"created"`
	Updated  int64     `json:"updated"`
	Findings []Finding `json:"findings"`
}

// ScanResultStore defines the operator for working with scan results.
type ScanResultStore interface {

	// Count returns the number of scan results from datastore.
	Count(ctx context.Context) (int64, error)

	// Find returns a scan result from datastore..
	Find(ctx context.Context, id int64) (*ScanResult, error)

	// Creates persists a scan result in the datastore.
	Create(ctx context.Context, s *ScanResult) error

	// List returns a list of scan result from the datastore.
	List(context.Context) ([]*ScanResult, error)
}

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
