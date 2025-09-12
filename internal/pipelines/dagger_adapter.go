// Package pipelines provides adapters to convert real Dagger types to interfaces.
package pipelines

import (
	"context"

	"dagger.io/dagger"
)

// DaggerAdapter adapts real Dagger types to our interfaces
type DaggerAdapter struct {
	client *dagger.Client
}

// GetRealClient returns the underlying real Dagger client for use with shared functions
func (a *DaggerAdapter) GetRealClient() *dagger.Client {
	return a.client
}

// NewDaggerAdapter creates a new adapter for a real Dagger client
func NewDaggerAdapter(client *dagger.Client) DaggerClient {
	return &DaggerAdapter{client: client}
}

func (a *DaggerAdapter) Host() DaggerHost {
	return &DaggerHostAdapter{host: a.client.Host()}
}

func (a *DaggerAdapter) Container() DaggerContainer {
	return &DaggerContainerAdapter{container: a.client.Container()}
}

func (a *DaggerAdapter) CacheVolume(name string) DaggerCacheVolume {
	return &DaggerCacheVolumeAdapter{cache: a.client.CacheVolume(name)}
}

// DaggerHostAdapter adapts dagger.Host to DaggerHost interface
type DaggerHostAdapter struct {
	host *dagger.Host
}

func (a *DaggerHostAdapter) Directory(path string, opts DaggerHostDirectoryOpts) DaggerDirectory {
	daggerOpts := dagger.HostDirectoryOpts{
		Exclude: opts.Exclude,
	}
	return &DaggerDirectoryAdapter{directory: a.host.Directory(path, daggerOpts)}
}

// DaggerDirectoryAdapter adapts dagger.Directory to DaggerDirectory interface
type DaggerDirectoryAdapter struct {
	directory *dagger.Directory
}

// GetRealDirectory returns the underlying real Dagger directory for use with shared functions
func (a *DaggerDirectoryAdapter) GetRealDirectory() *dagger.Directory {
	return a.directory
}

// DaggerContainerAdapter adapts dagger.Container to DaggerContainer interface
type DaggerContainerAdapter struct {
	container *dagger.Container
}

func (a *DaggerContainerAdapter) From(image string) DaggerContainer {
	return &DaggerContainerAdapter{container: a.container.From(image)}
}

func (a *DaggerContainerAdapter) WithMountedDirectory(path string, dir DaggerDirectory) DaggerContainer {
	// We need to extract the real directory from our adapter
	if adapter, ok := dir.(*DaggerDirectoryAdapter); ok {
		return &DaggerContainerAdapter{container: a.container.WithMountedDirectory(path, adapter.directory)}
	}
	return a
}

func (a *DaggerContainerAdapter) WithMountedCache(path string, cache DaggerCacheVolume) DaggerContainer {
	// We need to extract the real cache from our adapter
	if adapter, ok := cache.(*DaggerCacheVolumeAdapter); ok {
		return &DaggerContainerAdapter{container: a.container.WithMountedCache(path, adapter.cache)}
	}
	return a
}

func (a *DaggerContainerAdapter) WithWorkdir(path string) DaggerContainer {
	return &DaggerContainerAdapter{container: a.container.WithWorkdir(path)}
}

func (a *DaggerContainerAdapter) WithEnvVariable(name, value string) DaggerContainer {
	return &DaggerContainerAdapter{container: a.container.WithEnvVariable(name, value)}
}

func (a *DaggerContainerAdapter) WithExec(args []string, opts DaggerContainerWithExecOpts) DaggerContainer {
	daggerOpts := dagger.ContainerWithExecOpts{
		RedirectStdout: opts.RedirectStdout,
	}
	return &DaggerContainerAdapter{container: a.container.WithExec(args, daggerOpts)}
}

func (a *DaggerContainerAdapter) File(path string) DaggerFile {
	return &DaggerFileAdapter{file: a.container.File(path)}
}

// DaggerCacheVolumeAdapter adapts dagger.CacheVolume to DaggerCacheVolume interface
type DaggerCacheVolumeAdapter struct {
	cache *dagger.CacheVolume
}

// DaggerFileAdapter adapts dagger.File to DaggerFile interface
type DaggerFileAdapter struct {
	file *dagger.File
}

func (a *DaggerFileAdapter) Contents(ctx context.Context) (string, error) {
	return a.file.Contents(ctx)
}
