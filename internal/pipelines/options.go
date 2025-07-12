package pipelines

import "dagger.io/dagger"

// Package pipelines provides the core functionality for managing and executing different types of pipelines.

// Config holds configuration values used across all pipelines.
//
// Fields:
//   - Env: The environment in which the pipeline is running (e.g., dev, staging, prod).
//   - SkipPush: Indicates whether to skip the push stage.
//   - OnlyTest: Indicates whether to run only the test stage.
//   - OnlyBuild: Indicates whether to run only the build stage.
//   - Verbose: Enables verbose logging.
//   - GitRepo: The Git repository URL.
//   - GitRef: The Git reference (e.g., branch or tag).
//   - CIJobToken: The CI job token for authentication.
//   - RegistryURL: The Docker registry URL.
//   - RegistryToken: The Docker registry token.
//   - Version: The version of the application.
//   - BuildTag: The build tag for the application.
//   - CommitSHA: The commit SHA of the current build.
//   - BranchName: The name of the Git branch.
//   - Token: A generic token for authentication.
//   - Coverage: The test coverage percentage.
//   - GitUserEmail: The Git user email for commits.
//   - GitUserName: The Git user name for commits.
type Config struct {
	Env          string // The environment in which the pipeline is running.
	SkipPush     bool   // Indicates whether to skip the push stage.
	OnlyTest     bool   // Indicates whether to run only the test stage.
	OnlyBuild    bool   // Indicates whether to run only the build stage.
	Verbose      bool   // Enables verbose logging.
	GitRepo      string // The Git repository URL.
	GitRef       string // The Git reference (e.g., branch or tag).
	GitProtocol  string // The Git GitProtocol ssh or https.
	GitUserEmail string // The Git user email for commits.
	GitUserName  string // The Git user name for commits.

	RegistryURL    string  // The Docker registry URL.
	RegistryToken  string  // The Docker registry token.
	Version        string  // The version of the application.
	BuildTag       string  // The build tag for the application.
	CommitSHA      string  // The commit SHA of the current build.
	BranchName     string  // The name of the Git branch.
	Token          string  // A generic token for authentication.
	Coverage       float64 // The test coverage percentage.
	Image          *dagger.Container
	ImageRef       string
	ImageContainer *dagger.Container
	ImageName      string

	RegistryUser string // gitlab-ci-token (in CI) or your username
	RegistryPass string // CI_JOB_TOKEN (in CI) or your Personal Access Token
	Registry     string // Docker registry URL (e.g., registry.gitlab.com/my-org/my-project/service)
	ImageTag     string // Image tag (e.g., latest, sha, v1.2.3)

	GoVersion   string // Go version to use (e.g., 1.24.2)
	JavaVersion string // Java version to use (e.g., 17)

	SSHPrivateKey string // SSH private key for Git authentication
}

// Option represents a functional option for configuring Config.
//
// A functional option is a function that modifies a Config instance.
type Option func(*Config)

// NewConfig constructs a Config instance using functional options.
//
// Parameters:
//   - opts: A variadic list of Option functions to apply to the Config.
//
// Returns:
//   - A Config instance with the applied options.
func NewConfig(opts ...Option) Config {
	cfg := Config{
		Env:    "dev",  // Default environment.
		GitRef: "main", // Default branch.
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

// ====== Functional Options ======

// WithEnv sets the environment for the Config.
//
// Parameters:
//   - env: The environment to set (e.g., dev, staging, prod).
//
// Returns:
//   - An Option function to set the environment.
func WithEnv(env string) Option {
	return func(c *Config) { c.Env = env }
}

// WithSkipPush sets the SkipPush flag for the Config.
//
// Parameters:
//   - skip: A boolean indicating whether to skip the push stage.
//
// Returns:
//   - An Option function to set the SkipPush flag.
func WithSkipPush(skip bool) Option {
	return func(c *Config) { c.SkipPush = skip }
}

// WithOnlyTest sets the OnlyTest flag for the Config.
//
// Parameters:
//   - onlyTest: A boolean indicating whether to run only the test stage.
//
// Returns:
//   - An Option function to set the OnlyTest flag.
func WithOnlyTest(onlyTest bool) Option {
	return func(c *Config) { c.OnlyTest = onlyTest }
}

// WithOnlyBuild sets the OnlyBuild flag for the Config.
//
// Parameters:
//   - onlyBuild: A boolean indicating whether to run only the build stage.
//
// Returns:
//   - An Option function to set the OnlyBuild flag.
func WithOnlyBuild(onlyBuild bool) Option {
	return func(c *Config) { c.OnlyBuild = onlyBuild }
}

// WithVerbose sets the Verbose flag for the Config.
//
// Parameters:
//   - verbose: A boolean indicating whether to enable verbose logging.
//
// Returns:
//   - An Option function to set the Verbose flag.
func WithVerbose(verbose bool) Option {
	return func(c *Config) { c.Verbose = verbose }
}

// WithRegistry sets the Docker registry URL and token for the Config.
//
// Parameters:
//   - url: The Docker registry URL.
//   - token: The Docker registry token.
//
// Returns:
//   - An Option function to set the Docker registry URL and token.
func WithRegistry(url, token string) Option {
	return func(c *Config) {
		c.RegistryURL = url
		c.RegistryToken = token
	}
}

// WithRegistryUser sets the Docker registry user.
//
// Parameters:
//   - user: The Docker registry user.
//
// Returns:
//   - An Option function to set the Docker registry user.
func WithRegistryUser(user string) Option {
	return func(c *Config) {
		c.RegistryUser = user
	}
}

// WithBuildTag sets the build tag for the Config.
//
// Parameters:
//   - tag: The build tag to set.
//
// Returns:
//   - An Option function to set the build tag.
func WithBuildTag(tag string) Option {
	return func(c *Config) { c.BuildTag = tag }
}

// WithCommitSHA sets the commit SHA for the Config.
//
// Parameters:
//   - sha: The commit SHA to set.
//
// Returns:
//   - An Option function to set the commit SHA.
func WithCommitSHA(sha string) Option {
	return func(c *Config) { c.CommitSHA = sha }
}

// WithBranch sets the branch name for the Config.
//
// Parameters:
//   - branch: The branch name to set.
//
// Returns:
//   - An Option function to set the branch name.
func WithBranch(branch string) Option {
	return func(c *Config) { c.BranchName = branch }
}

// WithGitGitProtocol sets the Git protocol (e.g., "ssh" or "https") for the Config.
//
// Parameters:
//   - gitProtocol: The Git protocol to use ("ssh" or "https").
//
// Returns:
//   - An Option function that sets the GitProtocol field in Config.
func WithGitGitProtocol(gitProtocol string) Option {
	return func(c *Config) {
		c.GitProtocol = gitProtocol
	}
}

// WithGitRepo sets the Git repository URL and reference for the Config.
//
// Parameters:
//   - repo: The Git repository URL.
//   - ref: The Git reference (e.g., branch or tag).
//
// Returns:
//   - An Option function to set the Git repository URL and reference.
func WithGitRepo(repo, ref string) Option {
	return func(c *Config) {
		c.GitRepo = repo
		c.GitRef = ref
	}
}

// WithToken sets a generic token for the Config.
//
// Parameters:
//   - token: The token to set.
//
// Returns:
//   - An Option function to set the token.
func WithToken(token string) Option {
	return func(c *Config) { c.Token = token }
}

// WithCoverage sets the test coverage percentage for the Config.
//
// Parameters:
//   - percentage: The test coverage percentage to set.
//
// Returns:
//   - An Option function to set the test coverage percentage.
func WithCoverage(percentage float64) Option {
	return func(c *Config) { c.Coverage = percentage }
}

// WithGoVersion sets the Go version for the Config.
//
// Parameters:
//   - version: The Go version to set (e.g., "1.24.2").
//
// Returns:
//   - An Option function to set the Go version.
func WithGoVersion(version string) Option {
	return func(c *Config) { c.GoVersion = version }
}

// WithSSHPrivateKey sets the SSH private key for the Config.
//
// Parameters:
//   - key: The SSH private key to set.
//
// Returns:
//   - An Option function to set the SSH private key.
func WithSSHPrivateKey(key string) Option {
	return func(c *Config) { c.SSHPrivateKey = key }
}

// WithGitUserEmail sets the Git user email for the Config.
//
// Parameters:
//   - email: The Git user email to set.
//
// Returns:
//   - An Option function to set the Git user email.
func WithGitUserEmail(email string) Option {
	return func(c *Config) { c.GitUserEmail = email }
}

// WithGitUserName sets the Git user name for the Config.
//
// Parameters:
//   - name: The Git user name to set.
//
// Returns:
//   - An Option function to set the Git user name.
func WithGitUserName(name string) Option {
	return func(c *Config) { c.GitUserName = name }
}
