package db

import "errors"

// ErrOptimisticLock is returned by if the value is being modified by the database.
var ErrOptimisticLock = errors.New("Optimistic Lock Error")
