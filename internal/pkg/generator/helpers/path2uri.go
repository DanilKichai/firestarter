package helpers

import (
	"net/url"
	"path/filepath"
)

func Path2URI(pathParts ...string) string {
	path := "/"
	for _, pathPart := range pathParts {
		path = filepath.Join(path, pathPart)
	}

	uriParts := filepath.SplitList(path)

	uri := "file:///"
	for _, uriPart := range uriParts {
		uri, _ = url.JoinPath(uri, uriPart)
	}

	return uri
}
