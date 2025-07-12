// Package gokit provides the go-kit pipeline implementation.
package gokit

import (
	"context"
	"errors"
	"fmt"

	"dagger.io/dagger"

	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines/shared"
)

// GoKitPipeline represents a pipeline for the go-kit project.
//
// Fields:
//   - Client: The Dagger client used for container operations.
//   - Config: The configuration for the pipeline.
//   - Src: The source directory of the cloned repository.
//   - Cloner: The cloner used for cloning the repository.
type GoKitPipeline struct {
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

	var cloner shared.Cloner
	if cfg.GitProtocol == "ssh" {
		cloner = &shared.SSHCloner{}
	} else {
		cloner = &shared.HTTPSCloner{}
	}

	return &GoKitPipeline{
		Client: client,
		Config: cfg,
		Src:    src,
		Cloner: cloner,
	}
}

// Name returns the name of the pipeline.
//
// Returns:
//   - A string representing the name of the pipeline.
func (p *GoKitPipeline) Name() string {
	return "go-kit"
}

// Setup clones the go-kit repository and sets up the source directory.
//
// Parameters:
//   - ctx: The context for managing execution.
//
// Returns:
//   - An error if the cloning process fails, otherwise nil.
func (p *GoKitPipeline) Setup(ctx context.Context) error {
	if p.Cloner != nil {
		dir, err := p.Cloner.Clone(ctx, p.Client, shared.GitCloneOpts{})
		if err != nil {
			return err
		}
		p.Src = dir
	}
	return nil
}

// Test runs the tests for the go-kit project with coverage.
//
// Parameters:
//   - ctx: The context for managing execution.
//
// Returns:
//   - An error if the tests fail, otherwise nil.
func (p *GoKitPipeline) Test(ctx context.Context) error {
	if p.Src == nil {
		return errors.New("pipeline not set up: source directory is nil")
	}
	fmt.Println("🧪 running tests for go-kit...")
	return shared.RunTestsWithCoverage(ctx, p.Client, p.Src, p.Config.Coverage)
}

// Build is a placeholder for the build step of the pipeline.
//
// Parameters:
//   - ctx: The context for managing execution.
//
// Returns:
//   - Always returns nil as it is not implemented.
func (p *GoKitPipeline) Build(ctx context.Context) error {
	if p.Src == nil {
		return errors.New("pipeline not set up: source directory is nil")
	}
	builder := shared.NewGoBuilder(p.Client, p.Src, "1.21")
	outPath := "bin/app"
	_, err := builder.Build(ctx, outPath, outPath, map[string]string{"CGO_ENABLED": "0"})
	if err != nil {
		return fmt.Errorf("failed to build go-kit binary: %w", err)
	}
	fmt.Println("✅ go-kit binary built at", outPath)
	return nil
}

// Package is a placeholder for the packaging step of the pipeline.
//
// Parameters:
//   - ctx: The context for managing execution.
//
// Returns:
//   - Always returns nil as it is not implemented.
func (p *GoKitPipeline) Package(_ context.Context) error {
	return errors.New("not implemented")
}

// Tag generates a tag for the go-kit project.
//
// Parameters:
//   - ctx: The context for managing execution.
//
// Returns:
//   - An error if the tagging process fails, otherwise nil.
func (p *GoKitPipeline) Tag(ctx context.Context) error {
	fmt.Println("🧪 Generating tag for go-kit...")
	tag, err := shared.GenerateTag(ctx, p.Client, p.Src)
	if err != nil {
		return fmt.Errorf("❌ Error generating tag: %w", err)
	}
	fmt.Printf("✅ Tag generated: %s\n", tag)
	return nil
}

// Push is a placeholder for the push step of the pipeline.
//
// Parameters:
//   - ctx: The context for managing execution.
//
// Returns:
//   - Always returns nil as it is not implemented.
func (p *GoKitPipeline) Push(_ context.Context) error {
	return errors.New("not implemented")
}

// BeforeStep is a hook that executes custom logic before a specific pipeline step.
//
// Parameters:
//   - ctx: The context for managing execution.
//   - step: The name of the step.
//
// Returns:
//   - Always returns nil as it is not implemented.
func (p *GoKitPipeline) BeforeStep(_ context.Context, _ string) pipelines.HookFunc {
	return nil
}

// AfterStep is a hook that executes custom logic after a specific pipeline step.
//
// Parameters:
//   - ctx: The context for managing execution.
//   - step: The name of the step.
//
// Returns:
//   - Always returns nil as it is not implemented.
func (p *GoKitPipeline) AfterStep(_ context.Context, _ string) pipelines.HookFunc {
	return nil
}

func (p *GoKitPipeline) Cleanup(_ context.Context, _ *dagger.Client) error {
	return errors.New("not implemented")
}
