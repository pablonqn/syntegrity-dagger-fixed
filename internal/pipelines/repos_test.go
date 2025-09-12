package pipelines

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositories_Map(t *testing.T) {
	// Test that Repositories map is properly initialized
	assert.NotNil(t, Repositories)
	assert.Contains(t, Repositories, "go-kit")
	assert.Contains(t, Repositories, "docker-go")
}

func TestRepositories_GoKit(t *testing.T) {
	repo, exists := Repositories["go-kit"]
	assert.True(t, exists)
	assert.Equal(t, "https://gitlab.com/syntegrity/go-kit.git", repo.httpsURL)
	assert.Equal(t, "git@gitlab.com:syntegrity/go-kit.git", repo.sshURL)
}

func TestRepositories_DockerGo(t *testing.T) {
	repo, exists := Repositories["docker-go"]
	assert.True(t, exists)
	assert.Equal(t, "https://gitlab.com/syntegrity/docker-go.git", repo.httpsURL)
	assert.Equal(t, "git@gitlab.com:syntegrity/docker-go.git", repo.sshURL)
}

func TestGetRepoURL_HTTPS(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		protocol string
		expected string
	}{
		{
			name:     "go-kit https",
			repoName: "go-kit",
			protocol: "https",
			expected: "https://gitlab.com/syntegrity/go-kit.git",
		},
		{
			name:     "docker-go https",
			repoName: "docker-go",
			protocol: "https",
			expected: "https://gitlab.com/syntegrity/docker-go.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRepoURL(tt.repoName, tt.protocol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRepoURL_SSH(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		protocol string
		expected string
	}{
		{
			name:     "go-kit ssh",
			repoName: "go-kit",
			protocol: "ssh",
			expected: "git@gitlab.com:syntegrity/go-kit.git",
		},
		{
			name:     "docker-go ssh",
			repoName: "docker-go",
			protocol: "ssh",
			expected: "git@gitlab.com:syntegrity/docker-go.git",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRepoURL(tt.repoName, tt.protocol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRepoURL_NonexistentRepo(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		protocol string
		expected string
	}{
		{
			name:     "nonexistent repo https",
			repoName: "nonexistent",
			protocol: "https",
			expected: "",
		},
		{
			name:     "nonexistent repo ssh",
			repoName: "nonexistent",
			protocol: "ssh",
			expected: "",
		},
		{
			name:     "empty repo name",
			repoName: "",
			protocol: "https",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRepoURL(tt.repoName, tt.protocol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetRepoURL_InvalidProtocol(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		protocol string
		expected string
	}{
		{
			name:     "invalid protocol",
			repoName: "go-kit",
			protocol: "ftp",
			expected: "",
		},
		{
			name:     "empty protocol",
			repoName: "go-kit",
			protocol: "",
			expected: "",
		},
		{
			name:     "case sensitive protocol",
			repoName: "go-kit",
			protocol: "HTTPS",
			expected: "",
		},
		{
			name:     "case sensitive protocol ssh",
			repoName: "go-kit",
			protocol: "SSH",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRepoURL(tt.repoName, tt.protocol)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRepoInfo_Struct(t *testing.T) {
	// Test RepoInfo struct fields
	repo := RepoInfo{
		sshURL:   "git@gitlab.com:syntegrity/test.git",
		httpsURL: "https://gitlab.com/syntegrity/test.git",
	}

	assert.Equal(t, "git@gitlab.com:syntegrity/test.git", repo.sshURL)
	assert.Equal(t, "https://gitlab.com/syntegrity/test.git", repo.httpsURL)
}

func TestRepositories_Immutable(t *testing.T) {
	// Test that we can't accidentally modify the global Repositories map
	originalGoKit := Repositories["go-kit"]

	// Try to modify (this should not affect the original)
	repo := Repositories["go-kit"]
	repo.httpsURL = "modified"

	// Original should be unchanged
	assert.Equal(t, originalGoKit, Repositories["go-kit"])
	assert.NotEqual(t, repo, Repositories["go-kit"])
}

func TestGetRepoURL_AllRepositories(t *testing.T) {
	// Test all repositories in the map
	for repoName, repoInfo := range Repositories {
		t.Run(repoName+"_https", func(t *testing.T) {
			result := GetRepoURL(repoName, "https")
			assert.Equal(t, repoInfo.httpsURL, result)
		})

		t.Run(repoName+"_ssh", func(t *testing.T) {
			result := GetRepoURL(repoName, "ssh")
			assert.Equal(t, repoInfo.sshURL, result)
		})
	}
}

func TestRepositories_Consistency(t *testing.T) {
	// Test that all repositories have both SSH and HTTPS URLs
	for repoName, repoInfo := range Repositories {
		t.Run(repoName+"_consistency", func(t *testing.T) {
			assert.NotEmpty(t, repoInfo.httpsURL, "HTTPS URL should not be empty for %s", repoName)
			assert.NotEmpty(t, repoInfo.sshURL, "SSH URL should not be empty for %s", repoName)

			// Verify URLs are different
			assert.NotEqual(t, repoInfo.httpsURL, repoInfo.sshURL, "HTTPS and SSH URLs should be different for %s", repoName)

			// Verify HTTPS URL starts with https://
			assert.Contains(t, repoInfo.httpsURL, "https://", "HTTPS URL should start with https:// for %s", repoName)

			// Verify SSH URL starts with git@
			assert.Contains(t, repoInfo.sshURL, "git@", "SSH URL should start with git@ for %s", repoName)
		})
	}
}
