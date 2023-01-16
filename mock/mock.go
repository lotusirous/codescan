//go:build !oss
// +build !oss

package mock

//go:generate mockgen -package=mock -destination=mock_gen.go github.com/lotusirous/codescan/core RepositoryStore,ScanScheduler,ScanResultStore,ScanStore
