// Package pipelines provides the core functionality for managing and executing different types of pipelines.
package pipelines

import "context"

// HookFunc defines a function type that takes a context and returns an error.
// It is used as a hook to execute custom logic before or after a pipeline step.
type HookFunc func(ctx context.Context) error

// Pipeline defines an interface for a generic pipeline with methods
// to handle various stages of a pipeline process such as setup, testing,
// building, packaging, and pushing. It also includes hooks for executing
// custom logic before and after each step.
type Pipeline interface {
	// Test executes the test step of the pipeline.
	// It takes a context for managing execution and returns an error
	// if the tests fail.
	Test(ctx context.Context) error

	// Build compiles or builds the necessary components of the pipeline.
	// It takes a context for managing execution and returns an error
	// if the build process fails.
	Build(ctx context.Context) error

	// Package creates a distributable package for the pipeline.
	// It takes a context for managing execution and returns an error
	// if the packaging process fails.
	Package(ctx context.Context) error

	// Tag applies a tag to the pipeline's output.
	// It takes a context for managing execution and returns an error
	// if the tagging process fails.
	Tag(ctx context.Context) error

	// Name returns the name of the pipeline.
	// It returns a string representing the pipeline's name.
	Name() string

	// Setup performs the initial setup required for the pipeline.
	// It takes a context for managing execution and returns an error
	// if the setup fails.
	Setup(ctx context.Context) error

	// Push uploads or deploys the pipeline's output to a remote location.
	// It takes a context for managing execution and returns an error
	// if the push process fails.
	Push(ctx context.Context) error

	// BeforeStep is a hook that executes custom logic before a specific pipeline step.
	// It takes a context and the name of the step as arguments and returns a HookFunc.
	BeforeStep(ctx context.Context, step string) HookFunc

	// AfterStep is a hook that executes custom logic after a specific pipeline step.
	// It takes a context and the name of the step as arguments and returns a HookFunc.
	AfterStep(ctx context.Context, step string) HookFunc
}

// Setupper defines an interface for components that require a setup step.
// It includes a single method, Setup, which takes a context and returns an error
// if the setup process fails.
type Setupper interface {
	Setup(ctx context.Context) error
}

// Builder defines an interface for components that require a build step.
// It includes a single method, Build, which takes a context and returns an error
// if the build process fails.
type Builder interface {
	Build(ctx context.Context) error
}

// Tester defines an interface for components that require a test step.
// It includes a single method, Test, which takes a context and returns an error
// if the tests fail.
type Tester interface {
	Test(ctx context.Context) error
}

// Tagger defines an interface for components that require a tagging step.
// It includes a single method, Tag, which takes a context and returns an error
// if the tagging process fails.
type Tagger interface {
	Tag(ctx context.Context) error
}

// Packager defines an interface for components that require a packaging step.
// It includes a single method, Package, which takes a context and returns an error
// if the packaging process fails.
type Packager interface {
	Package(ctx context.Context) error
}

// Pusher defines an interface for components that require a push step.
// It includes a single method, Push, which takes a context and returns an error
// if the push process fails.
type Pusher interface {
	Push(ctx context.Context) error
}

// Package pipelines provides core pipeline interfaces and implementations.
