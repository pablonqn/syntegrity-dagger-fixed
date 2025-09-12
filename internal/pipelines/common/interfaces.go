// Package interfaces provides shared interfaces for pipeline components.
package interfaces

import "context"

// Pipeline defines the interface for a generic pipeline, providing methods
// for various stages of the pipeline process, such as setup, building, testing,
// packaging, tagging, and pushing. It also includes hooks for executing custom
// logic before and after each step.
type Pipeline interface {
	// Name returns the name of the pipeline.
	//
	// Returns:
	//   - A string representing the name of the pipeline.
	Name() string

	// Setup performs the initial setup required for the pipeline.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//
	// Returns:
	//   - An error if the setup process fails.
	Setup(ctx context.Context) error

	// Build compiles or builds the necessary components of the pipeline.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//
	// Returns:
	//   - An error if the build process fails.
	Build(ctx context.Context) error

	// Test executes the test step of the pipeline.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//
	// Returns:
	//   - An error if the tests fail.
	Test(ctx context.Context) error

	// Package creates a distributable package for the pipeline.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//
	// Returns:
	//   - An error if the packaging process fails.
	Package(ctx context.Context) error

	// Tag applies a tag to the pipeline's output.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//
	// Returns:
	//   - An error if the tagging process fails.
	Tag(ctx context.Context) error

	// Push uploads or deploys the pipeline's output to a remote location.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//
	// Returns:
	//   - An error if the push process fails.
	Push(ctx context.Context) error

	// BeforeStep is a hook that executes custom logic before a specific pipeline step.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//   - step: The name of the step before which the hook is executed.
	//
	// Returns:
	//   - A HookFunc to execute custom logic.
	BeforeStep(ctx context.Context, step string) HookFunc

	// AfterStep is a hook that executes custom logic after a specific pipeline step.
	//
	// Parameters:
	//   - ctx: The context for managing execution.
	//   - step: The name of the step after which the hook is executed.
	//
	// Returns:
	//   - A HookFunc to execute custom logic.
	AfterStep(ctx context.Context, step string) HookFunc
}

// HookFunc defines a function type that takes a context and returns an error.
// It is used as a hook to execute custom logic before or after a pipeline step.
type HookFunc func(ctx context.Context) error
