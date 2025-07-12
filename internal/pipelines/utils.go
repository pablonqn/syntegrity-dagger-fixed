package pipelines

import "fmt"

// GitlabNetrc generates a .netrc file content for GitLab authentication
func GitlabNetrc(user, token string) string {
	return fmt.Sprintf("machine gitlab.com login %s password %s\n", user, token)
}

type ReturnType = []int

func ExitOk() ReturnType {
	return []int{0}
}

func ExitOkOrTestFail() ReturnType {
	return []int{0, 1}
}
