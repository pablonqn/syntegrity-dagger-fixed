package shared

import (
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/require"
)

func TestHTTPSCloner_Clone(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cloner := NewHTTPSCloner()
	opts := GitCloneOpts{
		Repo:      "https://gitlab.com/test/repo.git",
		Branch:    "main",
		Name:      "test-repo",
		UserEmail: "test@example.com",
		UserName:  "Test User",
	}

	dir, err := cloner.Clone(ctx, client, opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to clone repository")
	require.Nil(t, dir)
}

func TestHTTPSCloner_Clone_Error(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cloner := NewHTTPSCloner()
	opts := GitCloneOpts{
		Repo:      "https://invalid-url.git",
		Branch:    "main",
		Name:      "test-repo",
		UserEmail: "test@example.com",
		UserName:  "Test User",
	}

	dir, err := cloner.Clone(ctx, client, opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to clone repository")
	require.Nil(t, dir)
}

func TestHTTPSCloner_Clone_InvalidURL(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cloner := NewHTTPSCloner()
	opts := GitCloneOpts{
		Repo:      "", // Empty repo URL
		Branch:    "main",
		Name:      "test-repo",
		UserEmail: "test@example.com",
		UserName:  "Test User",
	}

	dir, err := cloner.Clone(ctx, client, opts)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid repository URL: repo is empty")
	require.Nil(t, dir)
}
