package confx

import (
	"os"
	"path/filepath"
)

func getConfigPathFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	maxDepth := 5
	i := 0
	for i < maxDepth {
		_, err = os.Stat(filepath.Join(dir, "config.yml"))
		if err == nil {
			return dir, nil
		}
		dir = filepath.Join(dir, "..")
		i++
	}

	return "", err
}
