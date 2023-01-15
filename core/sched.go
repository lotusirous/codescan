package core

import "context"

// ScanScheduler schedules the scan, it delegates the scan for the worker.
type ScanScheduler interface {
	// Scan adds the repo to the queue, perform the scanning.
	ScanRepo(ctx context.Context, repo *Repository) (*Scan, error)

	// RestoreLastScan restores the last scan to the queue.
	RestoreLastScan(ctx context.Context) error

	// Starts runs the worker inside the manager.
	Start(ctx context.Context) error
}
