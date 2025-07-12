package pipelines

// RepoInfo represents information about a repository.
//
// Fields:
//   - URL: The URL of the repository.
type RepoInfo struct {
	sshURL   string // The URL of the repository.
	httpsURL string // The URL of the repository.
}

// Repositories is a map that associates repository names with their respective information.
//
// Keys:
//   - The name of the repository (e.g., "go-kit").
//
// Values:
//   - RepoInfo: The information about the repository, including its URL.
var Repositories = map[string]RepoInfo{
	"go-kit": {
		httpsURL: "https://gitlab.com/syntegrity/go-kit.git", // URL for the go-kit repository.
		sshURL:   "git@gitlab.com:syntegrity/go-kit.git",
	},
	"docker-go": {
		httpsURL: "https://gitlab.com/syntegrity/docker-go.git",
		sshURL:   "git@gitlab.com:syntegrity/docker-go.git",
	},
	// Add more repositories here as needed...
}

// GetRepoURL returns the repository URL by name and protocol ("ssh" or "https").
// Returns empty string if not found.
func GetRepoURL(name string, protocol string) string {
	repo, ok := Repositories[name]
	if !ok {
		return ""
	}
	switch protocol {
	case "ssh":
		return repo.sshURL
	case "https":
		return repo.httpsURL
	default:
		return ""
	}
}
