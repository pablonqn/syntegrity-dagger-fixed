// Package test provides testing utilities and implementations for pipelines.
package test

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"dagger.io/dagger"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines"
)

type GoTester struct {
	// Client is the Dagger client used for container operations.
	Client *dagger.Client
	// Src is the source directory of the cloned repository.
	Src *dagger.Directory
	// Config contains the configuration for the pipeline.
	Config pipelines.Config
	// MinCoverage is the minimum coverage percentage required for the tests.
	MinCoverage float64
}

// GoTester represents a Go testing environment.

var _ Testable = (*GoTester)(nil)

// RunTests executes the Go tests and checks the coverage.
func (g *GoTester) RunTests(ctx context.Context) error {
	goMod := g.Client.CacheVolume("go-mod-cache")
	goBuild := g.Client.CacheVolume("go-build-cache")

	base := g.Client.Container().
		From("golang:1.21-alpine").
		WithMountedDirectory("/app", g.Src).
		WithMountedCache("/go/pkg/mod", goMod).
		WithMountedCache("/root/.cache/go-build", goBuild).
		WithWorkdir("/app").
		WithEnvVariable("GOPATH", "/go").
		WithEnvVariable("GOCACHE", "/root/.cache/go-build")

	testContainer := base.WithExec([]string{"go", "test", "./...", "-v", "-coverprofile=coverage.out"}).
		WithExec([]string{"go", "tool", "cover", "-func=coverage.out"}, dagger.ContainerWithExecOpts{
			RedirectStdout: "/tmp/coverage.txt",
		})

	coverageFile := testContainer.File("/tmp/coverage.txt")
	content, err := coverageFile.Contents(ctx)
	if err != nil {
		return fmt.Errorf("❌ Error generating coverage report: %w", err)
	}

	fmt.Println("📄 Coverage report:\n" + content)

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "total:") {
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				raw := strings.TrimSuffix(fields[len(fields)-1], "%")
				coverage, err := strconv.ParseFloat(raw, 64)
				if err != nil {
					return fmt.Errorf("❌ Error parsing coverage: %w", err)
				}
				if coverage < g.MinCoverage {
					return fmt.Errorf("❌ Insufficient coverage: %.1f%% (minimum %.1f%%)", coverage, g.MinCoverage)
				}
				fmt.Printf("✅ Sufficient coverage: %.1f%% (minimum %.1f%%)\n", coverage, g.MinCoverage)
				return nil
			}
		}
	}

	return errors.New("❌ No line with total coverage found")
}
