package dockergo

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"dagger.io/dagger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines"
	"gitlab.com/syntegrity/syntegrity-infra/internal/pipelines/shared"
	"gitlab.com/syntegrity/syntegrity-infra/tests/mocks"
)

func TestDockerGoPipeline_Setup(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{
		GitProtocol:  "https",
		BranchName:   "main",
		GitUserEmail: "test@example.com",
		GitUserName:  "Test User",
	}

	mockCloner := &mocks.MockCloner{}
	mockCloner.MockClone = func(_ context.Context, client *dagger.Client, _ shared.GitCloneOpts) (*dagger.Directory, error) {
		dir := client.Directory()
		dir = dir.WithNewFile("go.mod", "module example.com/test\n\ngo 1.21\n")
		dir = dir.WithNewFile("main_test.go", `package main

import "testing"

func TestExample(t *testing.T) {
	if false {
		t.Errorf("should not fail")
	}
}`)
		return dir, nil
	}

	pipeline := New(client, cfg).(*DockerGoPipeline)
	pipeline.Cloner = mockCloner

	err = pipeline.Setup(ctx)
	require.NoError(t, err)
}

func TestDockerGoPipeline_Setup_Error(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{
		GitProtocol:  "https",
		BranchName:   "main",
		GitUserEmail: "test@example.com",
		GitUserName:  "Test User",
	}

	mockCloner := &mocks.MockCloner{}
	mockCloner.MockClone = func(_ context.Context, _ *dagger.Client, _ shared.GitCloneOpts) (*dagger.Directory, error) {
		return nil, errors.New("mock clone error")
	}

	pipeline := New(client, cfg).(*DockerGoPipeline)
	pipeline.Cloner = mockCloner

	err = pipeline.Setup(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mock clone error")
}

func TestDockerGoPipeline_Test(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	// Create a temporary directory with a simple Go test
	tmpDir := t.TempDir()

	// Create a simple test file
	testFile := `package main

import "testing"

func TestExample(t *testing.T) {
	t.Log("Example test passed")
}`
	if err := os.WriteFile(filepath.Join(tmpDir, "example_test.go"), []byte(testFile), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Create go.mod
	goMod := `module test

go 1.21`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatalf("failed to write go.mod: %v", err)
	}

	// Create pipeline
	pipeline := New(client, pipelines.Config{
		BranchName: "main",
		Registry:   "registry.gitlab.com/test",
		ImageTag:   "test",
	})

	// Set source directory
	src := client.Host().Directory(tmpDir)
	pipeline.(*DockerGoPipeline).Src = src

	// Run test
	err = pipeline.Test(ctx)
	require.NoError(t, err)
}

func TestDockerGoPipeline_Test_Error_NoSrc(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	// Create pipeline without source
	pipeline := New(client, pipelines.Config{
		BranchName: "main",
		Registry:   "registry.gitlab.com/test",
		ImageTag:   "test",
	})

	// Run test
	err = pipeline.Test(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "pipeline not set up: source directory is nil")
}

func TestDockerGoPipeline_Build(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{
		GitProtocol:  "https",
		BranchName:   "main",
		GitUserEmail: "test@example.com",
		GitUserName:  "Test User",
	}

	// Mock cloner para que el setup sea exitoso
	mockCloner := &mocks.MockCloner{}
	mockCloner.MockClone = func(_ context.Context, client *dagger.Client, _ shared.GitCloneOpts) (*dagger.Directory, error) {
		dir := client.Directory()
		dir = dir.WithNewFile("Dockerfile", "FROM alpine:latest\n")
		return dir, nil
	}

	pipeline := New(client, cfg).(*DockerGoPipeline)
	pipeline.Cloner = mockCloner
	_ = pipeline.Setup(ctx)

	err = pipeline.Build(ctx)
	require.NoError(t, err)
	require.NotNil(t, pipeline.Image)
}

func TestDockerGoPipeline_Build_Error_NoSrc(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{}
	pipeline := New(client, cfg).(*DockerGoPipeline)
	// No setup, Src es nil
	err = pipeline.Build(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "pipeline not set up: source directory is nil")
}

func TestDockerGoPipeline_Package(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{
		GitProtocol:  "https",
		BranchName:   "main",
		GitUserEmail: "test@example.com",
		GitUserName:  "Test User",
	}

	mockCloner := &mocks.MockCloner{}
	mockCloner.MockClone = func(_ context.Context, client *dagger.Client, _ shared.GitCloneOpts) (*dagger.Directory, error) {
		dir := client.Directory()
		dir = dir.WithNewFile("Dockerfile", "FROM alpine:latest\n")
		return dir, nil
	}

	pipeline := New(client, cfg).(*DockerGoPipeline)
	pipeline.Cloner = mockCloner
	_ = pipeline.Setup(ctx)
	_ = pipeline.Build(ctx)

	err = pipeline.Package(ctx)
	require.NoError(t, err)
}

func TestDockerGoPipeline_Tag_Error_NoImage(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{}
	pipeline := New(client, cfg).(*DockerGoPipeline)
	err = pipeline.Tag(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "image not built")
}

func TestDockerGoPipeline_Tag_Error_InvalidConfig(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{}
	pipeline := New(client, cfg).(*DockerGoPipeline)
	pipeline.Image = client.Container()
	err = pipeline.Tag(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid configuration")
}

func TestDockerGoPipeline_Push(t *testing.T) {
	// Setup
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	if err != nil {
		t.Fatalf("failed to connect to dagger: %v", err)
	}
	defer client.Close()

	cfg := pipelines.Config{
		GitProtocol:   "https",
		BranchName:    "main",
		GitUserEmail:  "test@example.com",
		GitUserName:   "Test User",
		RegistryURL:   "registry.example.com",
		ImageTag:      "test-tag",
		RegistryUser:  "test-user",
		RegistryToken: "test-token",
	}

	pipeline := New(client, cfg).(*DockerGoPipeline)
	pipeline.Src = client.Directory()

	// Test
	err = pipeline.Push(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no image built to push")
}

func TestDockerGoPipeline_Push_Error_NoRegistryUser(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	// Ensure GITLAB_CI is not set
	os.Unsetenv("GITLAB_CI")

	cfg := pipelines.Config{
		BranchName: "test-branch",
		Registry:   "registry.example.com",
		ImageTag:   "test-tag",
		// No RegistryUser set
	}
	pipeline := New(client, cfg).(*DockerGoPipeline)
	pipeline.Image = client.Container().From("alpine")

	err = pipeline.Push(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CI_REGISTRY_USER empty in local environment")
}

func TestDockerGoPipeline_Push_Error_NoLocalPassword(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	// Ensure GITLAB_CI is not set
	os.Unsetenv("GITLAB_CI")

	cfg := pipelines.Config{
		BranchName:   "test-branch",
		Registry:     "registry.example.com",
		ImageTag:     "test-tag",
		RegistryUser: "test-user",
		// No RegistryToken set
	}
	p := New(client, cfg).(*DockerGoPipeline)
	p.Image = client.Container().From("alpine")

	err = p.Push(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "CI_REGISTRY_USER empty in local environment")
}

func TestDockerGoPipeline_Push_Error_Publish(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	cfg := pipelines.Config{
		BranchName:    "test-branch",
		Registry:      "registry.example.com",
		ImageTag:      "test-tag",
		RegistryUser:  "test-user",
		RegistryToken: "test-token",
	}
	p := New(client, cfg).(*DockerGoPipeline)
	p.Image = client.Container().From("alpine")

	// Mock the container to fail on publish
	p.Image = p.Image.WithExec([]string{"exit", "1"})

	err = p.Push(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error when pushing the image")
}

func TestDockerGoPipeline_Push_Success_GitLabCI(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	t.Setenv("GITLAB_CI", "true")
	t.Setenv("CI_JOB_TOKEN", "test-token")
	defer func() {
		os.Unsetenv("GITLAB_CI")
		os.Unsetenv("CI_JOB_TOKEN")
	}()

	cfg := pipelines.Config{
		BranchName: "test-branch",
		Registry:   "registry.example.com",
		ImageTag:   "test-tag",
	}
	p := New(client, cfg).(*DockerGoPipeline)
	p.Image = client.Container().From("alpine")

	err = p.Push(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error when pushing the image")
}

func TestDockerGoPipeline_Push_Success_Local(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	// Ensure GITLAB_CI is not set
	os.Unsetenv("GITLAB_CI")

	cfg := pipelines.Config{
		BranchName:    "test-branch",
		Registry:      "registry.example.com",
		ImageTag:      "test-tag",
		RegistryUser:  "test-user",
		RegistryToken: "test-token",
	}
	p := New(client, cfg).(*DockerGoPipeline)
	p.Image = client.Container().From("alpine")

	err = p.Push(ctx)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error when pushing the image")
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  pipelines.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: pipelines.Config{
				BranchName: "main",
				Registry:   "registry.gitlab.com/test",
				ImageTag:   "v1.0.0",
			},
			wantErr: false,
		},
		{
			name: "missing branch name",
			config: pipelines.Config{
				Registry: "registry.gitlab.com/test",
				ImageTag: "v1.0.0",
			},
			wantErr: true,
		},
		{
			name: "missing registry",
			config: pipelines.Config{
				BranchName: "main",
				ImageTag:   "v1.0.0",
			},
			wantErr: true,
		},
		{
			name: "missing image tag",
			config: pipelines.Config{
				BranchName: "main",
				Registry:   "registry.gitlab.com/test",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  pipelines.Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDockerGoPipeline_BeforeStep(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	cfg := pipelines.Config{
		BranchName: "test-branch",
		Registry:   "registry.example.com",
		ImageTag:   "test-tag",
	}
	pipeline := New(client, cfg)
	hook := pipeline.BeforeStep(ctx, "test-step")
	assert.Nil(t, hook)
}

func TestDockerGoPipeline_AfterStep(t *testing.T) {
	ctx := t.Context()
	client, err := dagger.Connect(ctx)
	require.NoError(t, err)
	defer client.Close()

	cfg := pipelines.Config{
		BranchName: "test-branch",
		Registry:   "registry.example.com",
		ImageTag:   "test-tag",
	}
	pipeline := New(client, cfg)
	hook := pipeline.AfterStep(ctx, "test-step")
	assert.Nil(t, hook)
}
