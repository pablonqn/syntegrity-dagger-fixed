// Package shared contiene utilidades y helpers para operaciones comunes en los pipelines.
// shared/cloner.go
package shared

import (
	"context"
	"errors"

	"dagger.io/dagger"
)

// Cloner define la interfaz para clonar repositorios
type Cloner interface {
	Clone(ctx context.Context, client *dagger.Client, opts GitCloneOpts) (*dagger.Directory, error)
}

// GitCloneOpts contiene las opciones para clonar un repositorio
type GitCloneOpts struct {
	Repo      string // URL del repositorio
	Branch    string // Rama a clonar
	Name      string // Nombre del directorio destino
	UserEmail string // Email del usuario Git
	UserName  string // Nombre del usuario Git
	SSHKey    string // Clave SSH (opcional)
}

// RepoCloner interface allows multiple auth modes.
type RepoCloner interface {
	Clone(ctx context.Context, client *dagger.Client, opts GitCloneOpts) (*dagger.Directory, error)
}

// CloneRepo es una función helper que usa el clonador por defecto
func CloneRepo(ctx context.Context, client *dagger.Client, opts GitCloneOpts, protocol string) (*dagger.Directory, error) {
	if err := validateOpts(opts); err != nil {
		return nil, err
	}

	var cloner Cloner
	if protocol == "ssh" {
		cloner = &SSHCloner{}
	} else {
		cloner = &HTTPSCloner{}
	}
	return cloner.Clone(ctx, client, opts)
}

func validateOpts(opts GitCloneOpts) error {
	if opts.Repo == "" || opts.Branch == "" || opts.Name == "" {
		return errors.New("❌ invalid GitCloneOpts: Repo, Branch and Name must be set")
	}
	return nil
}
