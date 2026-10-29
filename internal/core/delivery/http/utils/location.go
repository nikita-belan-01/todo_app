package utils

import (
	"net/http"
	"net/url"
	"strings"
)

func CreatedLocation(r *http.Request, segments ...string) string {
	path, _, _ := strings.Cut(r.RequestURI, "?")

	var location strings.Builder
	location.WriteString(strings.TrimSuffix(path, "/"))

	for _, segment := range segments {
		location.WriteByte('/')
		location.WriteString(url.PathEscape(segment))
	}

	return location.String()
}
