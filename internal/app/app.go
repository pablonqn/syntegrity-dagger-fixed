package app

import (
	"context"
	"fmt"
	"sync"

	kitlog "github.com/getsyntegrity/go-kit-logger/pkg/logger"
	"gitlab.com/syntegrity/syntegrity-infra/internal/interfaces"
)

var (
	globalContainer *Container
	containerOnce   sync.Once
	containerMutex  sync.RWMutex
)

// GetContainer returns the global container instance.
func GetContainer() *Container {
	containerMutex.RLock()
	defer containerMutex.RUnlock()

	if globalContainer == nil {
		panic("Container not initialized. Call app.Initialize() first")
	}
	return globalContainer
}

// Initialize sets up the global container instance.
func Initialize(ctx context.Context, cfg interfaces.Configuration) error {
	var initErr error

	containerOnce.Do(func() {
		container := NewContainer(ctx, cfg)

		// Initialize and set global logger
		logger := container.CreateLogger()
		kitlog.SetGlobal(logger)

		if err := container.Start(ctx); err != nil {
			initErr = fmt.Errorf("failed to start container: %w", err)
			return
		}

		containerMutex.Lock()
		globalContainer = container
		containerMutex.Unlock()
	})

	return initErr
}

// Reset clears the global container (useful for testing).
func Reset() {
	containerMutex.Lock()
	if globalContainer != nil {
		_ = globalContainer.Stop(context.Background())
	}
	globalContainer = nil
	containerOnce = sync.Once{} // Reset the sync.Once to allow re-initialization
	containerMutex.Unlock()
}

// App represents the main application instance.
// It manages the application lifecycle and provides access to core functionality.
//
// Features:
// - Application lifecycle management (start/stop)
// - Container management
// - Graceful shutdown
// - Component health monitoring
//
// Example usage:
//
//	container := app.GetContainer()
//	application := app.NewApp(container)
//	if err := application.Start(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer application.Stop(ctx)
type App struct {
	container *Container
}

// NewApp creates a new application instance.
func NewApp(container *Container) *App {
	return &App{container: container}
}

// Start starts the application.
func (a *App) Start(ctx context.Context) error {
	logger, err := a.container.GetLogger()
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	logger.Info("Starting Syntegrity Dagger application...")
	return a.container.Start(ctx)
}

// Stop stops the application.
func (a *App) Stop(ctx context.Context) error {
	logger, err := a.container.GetLogger()
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	logger.Info("Stopping Syntegrity Dagger application...")
	return a.container.Stop(ctx)
}

// GetContainer returns the application's container.
func (a *App) GetContainer() *Container {
	return a.container
}

// RunPipeline executes a pipeline with the given name.
func (a *App) RunPipeline(ctx context.Context, pipelineName string) error {
	logger, err := a.container.GetLogger()
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	logger.Info("Running pipeline", "name", pipelineName)

	pipeline, err := a.container.GetPipeline(pipelineName)
	if err != nil {
		return fmt.Errorf("failed to get pipeline %s: %w", pipelineName, err)
	}

	// Execute pipeline steps
	steps := pipeline.GetAvailableSteps()

	for _, step := range steps {
		logger.Info("Executing pipeline step", "step", step)

		stepErr := pipeline.ExecuteStep(ctx, step)
		if stepErr != nil {
			return fmt.Errorf("pipeline step %s failed: %w", step, stepErr)
		}

		logger.Info("Pipeline step completed", "step", step)
	}

	logger.Info("Pipeline completed successfully", "name", pipelineName)
	return nil
}

// RunPipelineStep executes a single pipeline step.
func (a *App) RunPipelineStep(ctx context.Context, pipelineName, step string) error {
	logger, err := a.container.GetLogger()
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	logger.Info("Running pipeline step", "pipeline", pipelineName, "step", step)

	pipeline, err := a.container.GetPipeline(pipelineName)
	if err != nil {
		return fmt.Errorf("failed to get pipeline %s: %w", pipelineName, err)
	}

	// Execute single step
	stepErr := pipeline.ExecuteStep(ctx, step)
	if stepErr != nil {
		return fmt.Errorf("pipeline step %s failed: %w", step, stepErr)
	}

	logger.Info("Pipeline step completed successfully", "pipeline", pipelineName, "step", step)
	return nil
}

// ListPipelines returns a list of available pipelines.
func (a *App) ListPipelines() ([]string, error) {
	registry, err := a.container.GetPipelineRegistry()
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline registry: %w", err)
	}

	return registry.List(), nil
}

// GetPipelineInfo returns information about a specific pipeline.
func (a *App) GetPipelineInfo(pipelineName string) (map[string]any, error) {
	pipeline, err := a.container.GetPipeline(pipelineName)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline %s: %w", pipelineName, err)
	}

	info := map[string]any{
		"name": pipeline.Name(),
		"steps": []string{
			"setup", "build", "test", "tag", "package", "push", "lint", "security", "release",
		},
	}

	return info, nil
}
