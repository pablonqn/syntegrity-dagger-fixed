// Package pipelines provides interfaces for Dagger types to enable mocking in tests.
package pipelines

import (
	"context"

	"dagger.io/dagger"
)

// DaggerClient interface abstracts the dagger.Client to enable mocking
type DaggerClient interface {
	Host() DaggerHost
	Container() DaggerContainer
	CacheVolume(string) DaggerCacheVolume
	// GetRealClient returns the underlying real Dagger client (only for adapters)
	GetRealClient() *dagger.Client
}

// DaggerHost interface abstracts the dagger.Host to enable mocking
type DaggerHost interface {
	Directory(string, DaggerHostDirectoryOpts) DaggerDirectory
}

// DaggerHostDirectoryOpts represents options for host directory operations
type DaggerHostDirectoryOpts struct {
	Exclude []string
}

// DaggerDirectory interface abstracts the dagger.Directory to enable mocking
type DaggerDirectory interface {
	// GetRealDirectory returns the underlying real Dagger directory (only for adapters)
	GetRealDirectory() *dagger.Directory
}

// DaggerContainer interface abstracts the dagger.Container to enable mocking
type DaggerContainer interface {
	From(string) DaggerContainer
	WithMountedDirectory(string, DaggerDirectory) DaggerContainer
	WithMountedCache(string, DaggerCacheVolume) DaggerContainer
	WithWorkdir(string) DaggerContainer
	WithEnvVariable(string, string) DaggerContainer
	WithExec([]string, DaggerContainerWithExecOpts) DaggerContainer
	File(string) DaggerFile
}

// DaggerContainerWithExecOpts represents options for container exec operations
type DaggerContainerWithExecOpts struct {
	RedirectStdout string
}

// DaggerCacheVolume interface abstracts the dagger.CacheVolume to enable mocking
type DaggerCacheVolume interface {
	// Add methods that are actually used in the pipelines
}

// DaggerFile interface abstracts the dagger.File to enable mocking
type DaggerFile interface {
	Contents(context.Context) (string, error)
}
