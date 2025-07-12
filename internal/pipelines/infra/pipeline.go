// Package infra provides infrastructure-related pipeline implementations.
package infra

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines/shared"
)

// Package infra implementa el pipeline para la infraestructura de Syntegrity.

// SyntegrityInfraPipeline represents a pipeline for the syntegrity-infra pipeline
//
// Fields:
//   - Client: The Dagger client used for container operations.
//   - Config: The configuration for the pipeline.
//   - Src: The source directory of the cloned repository.
//   - Cloner: The cloner used for cloning repositories.
type SyntegrityInfraPipeline struct {
	Client *dagger.Client
	Config pipelines.Config
	Src    *dagger.Directory
	Cloner shared.Cloner
}

// New creates a new instance of GoKitPipeline.
//
// Parameters:
//   - client: The Dagger client used for container operations.
//   - cfg: The configuration for the pipeline.
//
// Returns:
//   - A new instance of GoKitPipeline.
func New(client *dagger.Client, cfg pipelines.Config) pipelines.Pipeline {
	src := client.Host().Directory(".", dagger.HostDirectoryOpts{
		Exclude: []string{"**/node_modules", "**/.git", "**/.dagger-cache"},
	})
	return &SyntegrityInfraPipeline{
		Client: client,
		Config: cfg,
		Src:    src,
	}
}

func (s *SyntegrityInfraPipeline) Test(_ context.Context) error {
	return nil
}

func (s *SyntegrityInfraPipeline) Build(_ context.Context) error {
	return nil
}

func (s *SyntegrityInfraPipeline) Package(_ context.Context) error {
	return nil
}

func (s *SyntegrityInfraPipeline) Tag(_ context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (s *SyntegrityInfraPipeline) Name() string {
	return "syntegrity-infra"
}

func (s *SyntegrityInfraPipeline) Setup(ctx context.Context) error {
	fmt.Println("🧪 Running go-kit tests syntegrity-infra.....")
	if s.Cloner != nil {
		dir, err := s.Cloner.Clone(ctx, s.Client, shared.GitCloneOpts{})
		if err != nil {
			return err
		}
		s.Src = dir
	}
	return shared.RunTestsWithCoverage(ctx, s.Client, s.Src, s.Config.Coverage)
}

func (s *SyntegrityInfraPipeline) Push(_ context.Context) error {
	// TODO implement me
	panic("implement me")
}

func (s *SyntegrityInfraPipeline) BeforeStep(_ context.Context, _ string) pipelines.HookFunc {
	// TODO implement me
	panic("implement me")
}

func (s *SyntegrityInfraPipeline) AfterStep(_ context.Context, _ string) pipelines.HookFunc {
	// TODO implement me
	panic("implement me")
}
