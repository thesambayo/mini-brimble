package api

import (
	"path"
	"strings"
)

func GetAppNameFromRepo(repoURL string) string {
	// 1. Handle SSH format (git@github.com:user/repo.git)
	// We convert the colon to a slash to make it look like a standard path
	normalized := strings.ReplaceAll(repoURL, ":", "/")

	// 2. Use path.Base to get the last element
	// path.Base works on any string with forward slashes
	name := path.Base(normalized)

	// 3. Clean up the .git suffix
	name = strings.TrimSuffix(name, ".git")

	// 4. Safety check: handle empty or invalid inputs
	if name == "." || name == "/" {
		return ""
	}

	return name
}
