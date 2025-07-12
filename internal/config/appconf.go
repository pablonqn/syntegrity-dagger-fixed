// Package config contiene la configuración de la aplicación Syntegrity.
package config

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

const (
	envPrefix = "GETFIVVY_INFRA_SVC_"
)

var k = koanf.New("_")

// AppConf holds the full application configuration
type AppConf struct {
	// Log configuration
	LogLevel            string
	LogEncoding         string
	LogSamplingEnable   bool
	InfoLogSamplingRate int

	// GoVersion to be used in pipelines or builds
	GoVersion string

	// GitLab container registry info
	RegistryBaseURL string
	RegistryUser    string
	RegistryPass    string
}

// Init reads and parses multiple configuration sources and produces an instance of *AppConf
func Init() (*AppConf, error) {
	if err := loadDotEnv(); err != nil {
		return nil, err
	}

	err := k.Load(env.Provider(envPrefix, "_", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, envPrefix))
	}), nil)
	if err != nil {
		return nil, err
	}

	c := defaultAppConfig()
	err = k.Unmarshal("", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// defaultAppConfig creates new config with sane defaults to reduce the amount of local changes
func defaultAppConfig() *AppConf {
	defaultCfg := &AppConf{
		LogLevel:            "debug",
		LogEncoding:         "json",
		LogSamplingEnable:   false,
		InfoLogSamplingRate: 50,
		GoVersion:           "1.24.3",
		RegistryBaseURL:     "registry.gitlab.com/syntegrity",
		RegistryUser:        "pablo.martin.gore@gmail.com",
	}

	return defaultCfg
}

// VcsRevision returns the VCS revision from the build info
func (c *AppConf) VcsRevision() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	revision := ""
	for _, setting := range buildInfo.Settings {
		if setting.Key == "vcs.revision" {
			revision = setting.Value
			break
		}
	}

	return revision
}

// loadDotEnv loads the .env file if it exists
func loadDotEnv() error {
	p, err := os.Executable()
	if err != nil {
		return nil
	}

	envFile := filepath.Join(filepath.Dir(p), ".env")
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return nil
	}

	return godotenv.Load(envFile)
}
