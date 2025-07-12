// Package test provides testing utilities and implementations for pipelines.
package test

import (
	"context"
	"fmt"

	"dagger.io/dagger"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines"
)

// Testable is an interface that defines a contract for running tests.
type Testable interface {
	RunTests(ctx context.Context) error
}

// New returns a Testable implementation based on the configured language.
// Extend this switch to support other languages.
func New(client *dagger.Client, src *dagger.Directory, cfg pipelines.Config, language string) Testable {
	switch language {
	case "go":
		return &GoTester{
			Client:      client,
			Src:         src,
			Config:      cfg,
			MinCoverage: cfg.Coverage,
		}
	default:
		return &noopTester{}
	}
}

// noopTester implements Testable but does nothing — fallback or placeholder.
type noopTester struct{}

func (n *noopTester) RunTests(_ context.Context) error {
	fmt.Println("⚠️ No valid tester available for selected language.")
	return nil
}
