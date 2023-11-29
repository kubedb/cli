package filepath

import (
	sj "github.com/cyphar/filepath-securejoin"
)

func SecureJoin(root string, unsafePath ...string) (string, error) {
	result := root
	var err error
	for _, dir := range unsafePath {
		result, err = sj.SecureJoin(result, dir)
		if err != nil {
			return "", err
		}
	}
	return result, nil
}
