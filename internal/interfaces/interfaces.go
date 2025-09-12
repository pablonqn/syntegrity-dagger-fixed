package interfaces

import (
	"context"
	"time"

	"dagger.io/dagger"
)

// Configuration defines the application configuration interface.
type Configuration interface {
	Load() error
	LoadWithDefaults(defaults map[string]any) error
	Validate() error
	GetConfigSummary() string

	// Core configuration methods
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat(key string) float64
	Get(key string) any
	Set(key string, value any)
	Has(key string) bool
	All() map[string]any

	// Pipeline-specific configuration methods
	Pipeline() PipelineConfig
	Registry() RegistryConfig
	Security() SecurityConfig
	Logging() LoggingConfig
	Environment() string
}

// PipelineConfig defines pipeline configuration.
type PipelineConfig struct {
	Name        string  `yaml:"name" json:"name"`
	Environment string  `yaml:"environment" json:"environment"`
	Coverage    float64 `yaml:"coverage" json:"coverage"`
	SkipPush    bool    `yaml:"skip_push" json:"skip_push"`
	OnlyBuild   bool    `yaml:"only_build" json:"only_build"`
	OnlyTest    bool    `yaml:"only_test" json:"only_test"`
	Verbose     bool    `yaml:"verbose" json:"verbose"`
	GoVersion   string  `yaml:"go_version" json:"go_version"`
	JavaVersion string  `yaml:"java_version" json:"java_version"`
}

// RegistryConfig defines container registry configuration.
type RegistryConfig struct {
	BaseURL string `yaml:"base_url" json:"base_url"`
	User    string `yaml:"user" json:"user"`
	Pass    string `yaml:"pass" json:"pass"`
	Image   string `yaml:"image" json:"image"`
	Tag     string `yaml:"tag" json:"tag"`
}

// SecurityConfig defines security configuration.
type SecurityConfig struct {
	EnableVulnCheck bool     `yaml:"enable_vuln_check" json:"enable_vuln_check"`
	EnableLinting   bool     `yaml:"enable_linting" json:"enable_linting"`
	LintTimeout     string   `yaml:"lint_timeout" json:"lint_timeout"`
	ExcludePatterns []string `yaml:"exclude_patterns" json:"exclude_patterns"`
}

// LoggingConfig defines logging configuration.
type LoggingConfig struct {
	Level            string        `yaml:"level" json:"level"`
	Format           string        `yaml:"format" json:"format"`
	SamplingEnable   bool          `yaml:"sampling_enable" json:"sampling_enable"`
	SamplingRate     float64       `yaml:"sampling_rate" json:"sampling_rate"`
	SamplingInterval time.Duration `yaml:"sampling_interval" json:"sampling_interval"`
}

// GitConfig defines Git configuration.
type GitConfig struct {
	Repo      string `yaml:"repo" json:"repo"`
	Ref       string `yaml:"ref" json:"ref"`
	Protocol  string `yaml:"protocol" json:"protocol"`
	UserEmail string `yaml:"user_email" json:"user_email"`
	UserName  string `yaml:"user_name" json:"user_name"`
	SSHKey    string `yaml:"ssh_key" json:"ssh_key"`
}

// ReleaseConfig defines release configuration.
type ReleaseConfig struct {
	Enabled        bool     `yaml:"enabled" json:"enabled"`
	UseGoreleaser  bool     `yaml:"use_goreleaser" json:"use_goreleaser"`
	BuildTargets   []string `yaml:"build_targets" json:"build_targets"`
	ArchiveFormats []string `yaml:"archive_formats" json:"archive_formats"`
	Checksum       bool     `yaml:"checksum" json:"checksum"`
	Sign           bool     `yaml:"sign" json:"sign"`
}

// Container defines the interface for the dependency injection container.
type Container interface {
	// Core container methods
	Register(name string, factory func() (any, error))
	Get(name string) (any, error)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Validate() error

	// Pipeline-specific providers
	PipelineProvider
	RegistryProvider
	SecurityProvider
	LoggingProvider
}

// PipelineProvider defines the interface for pipeline access.
type PipelineProvider interface {
	GetPipelineRegistry() (PipelineRegistry, error)
	GetPipeline(name string) (Pipeline, error)
	GetDaggerClient() (*dagger.Client, error)
}

// RegistryProvider defines the interface for registry access.
type RegistryProvider interface {
	GetRegistryConfig() (RegistryConfig, error)
	GetRegistryAuth() (string, string, error)
}

// SecurityProvider defines the interface for security tools access.
type SecurityProvider interface {
	GetVulnChecker() (VulnChecker, error)
	GetLinter() (Linter, error)
}

// LoggingProvider defines the interface for logging access.
type LoggingProvider interface {
	GetLogger() (Logger, error)
}

// Pipeline defines the core pipeline interface with dynamic step execution.
type Pipeline interface {
	Name() string
	GetAvailableSteps() []string
	ExecuteStep(ctx context.Context, stepName string) error
	BeforeStep(ctx context.Context, stepName string) HookFunc
	AfterStep(ctx context.Context, stepName string) HookFunc
	GetStepConfig(stepName string) StepConfig
	ValidateStep(stepName string) error
}

// StepConfig defines configuration for a pipeline step.
type StepConfig struct {
	Name        string            `yaml:"name" json:"name"`
	Description string            `yaml:"description" json:"description"`
	Required    bool              `yaml:"required" json:"required"`
	Parallel    bool              `yaml:"parallel" json:"parallel"`
	Timeout     time.Duration     `yaml:"timeout" json:"timeout"`
	Retries     int               `yaml:"retries" json:"retries"`
	DependsOn   []string          `yaml:"depends_on" json:"depends_on"`
	Conditions  map[string]string `yaml:"conditions" json:"conditions"`
	Metadata    map[string]any    `yaml:"metadata" json:"metadata"`
}

// StepExecutor defines the interface for executing pipeline steps.
type StepExecutor interface {
	Execute(ctx context.Context, stepName string, config StepConfig) error
	GetStepResult(stepName string) (StepResult, error)
	GetStepLogs(stepName string) ([]string, error)
}

// StepResult contains the result of a step execution.
type StepResult struct {
	StepName  string         `json:"step_name"`
	Success   bool           `json:"success"`
	Duration  time.Duration  `json:"duration"`
	Error     error          `json:"error,omitempty"`
	Output    string         `json:"output,omitempty"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Artifacts []Artifact     `json:"artifacts,omitempty"`
}

// Artifact represents a file or artifact produced by a step.
type Artifact struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
	Checksum    string `json:"checksum,omitempty"`
	Description string `json:"description,omitempty"`
}

// PipelineRegistry defines the interface for pipeline registration and retrieval.
type PipelineRegistry interface {
	Register(name string, factory func(*dagger.Client, Configuration) Pipeline)
	Get(name string, client *dagger.Client, cfg Configuration) (Pipeline, error)
	List() []string
}

// VulnChecker defines the interface for vulnerability checking.
type VulnChecker interface {
	Check(ctx context.Context, src *dagger.Directory) error
	GetReport(ctx context.Context) (string, error)
}

// Linter defines the interface for code linting.
type Linter interface {
	Lint(ctx context.Context, src *dagger.Directory) error
	GetReport(ctx context.Context) (string, error)
}

// Logger defines the interface for logging.
type Logger interface {
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	Fatal(msg string, fields ...any)
	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}

// HookFunc defines a function type for pipeline hooks.
type HookFunc func(ctx context.Context) error

// StepHandler defines the interface for handling pipeline steps.
type StepHandler interface {
	CanHandle(stepName string) bool
	Execute(ctx context.Context, stepName string, config StepConfig) error
	GetStepInfo(stepName string) StepConfig
	Validate(stepName string, config StepConfig) error
}

// HookManager defines the interface for managing pipeline hooks.
type HookManager interface {
	RegisterHook(stepName string, hookType HookType, hook HookFunc) error
	GetHooks(stepName string, hookType HookType) []HookFunc
	ExecuteHooks(ctx context.Context, stepName string, hookType HookType) error
	RemoveHook(stepName string, hookType HookType, hook HookFunc) error
}

// HookType defines the type of hook.
type HookType string

const (
	HookTypeBefore  HookType = "before"
	HookTypeAfter   HookType = "after"
	HookTypeError   HookType = "error"
	HookTypeSuccess HookType = "success"
)

// StepRegistry defines the interface for registering and managing steps.
type StepRegistry interface {
	RegisterStep(stepName string, handler StepHandler) error
	GetStepHandler(stepName string) (StepHandler, error)
	ListSteps() []string
	GetStepConfig(stepName string) (StepConfig, error)
	ValidateStep(stepName string) error
	ExecuteStep(ctx context.Context, stepName string) error
	GetExecutionOrder() ([]string, error)
}

// PipelineExecutor defines the interface for executing pipelines with dynamic steps.
type PipelineExecutor interface {
	ExecutePipeline(ctx context.Context, pipelineName string, steps []string) error
	ExecuteStep(ctx context.Context, pipelineName string, stepName string) error
	GetPipelineStatus(pipelineName string) (PipelineStatus, error)
	GetStepStatus(pipelineName string, stepName string) (StepResult, error)
	CancelPipeline(pipelineName string) error
	GetPipelineLogs(pipelineName string) ([]string, error)
}

// PipelineStatus represents the current status of a pipeline.
type PipelineStatus struct {
	PipelineName string                `json:"pipeline_name"`
	Status       string                `json:"status"` // running, completed, failed, cancelled
	StartTime    time.Time             `json:"start_time"`
	EndTime      *time.Time            `json:"end_time,omitempty"`
	Duration     time.Duration         `json:"duration"`
	Steps        map[string]StepResult `json:"steps"`
	Metadata     map[string]any        `json:"metadata,omitempty"`
}

// ContainerError represents an error that occurred in the container.
type ContainerError struct {
	Component string
	Operation string
	Cause     error
}

func (e ContainerError) Error() string {
	return "container error in " + e.Component + " during " + e.Operation + ": " + e.Cause.Error()
}

// PipelineError represents a pipeline-specific error.
type PipelineError struct {
	Pipeline string
	Step     string
	Message  string
	Cause    error
}

func (e PipelineError) Error() string {
	return "pipeline error in " + e.Pipeline + " at step " + e.Step + ": " + e.Message + ": " + e.Cause.Error()
}

// ConfigurationError represents a configuration error.
type ConfigurationError struct {
	Key   string
	Value any
	Cause error
}

func (e ConfigurationError) Error() string {
	return "configuration error for key " + e.Key + " with value " + e.Value.(string) + ": " + e.Cause.Error()
}
