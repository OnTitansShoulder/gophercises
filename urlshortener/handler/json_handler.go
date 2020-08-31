package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var redirects []RedirectEntry

	err := json.Unmarshal(jsonBytes, &redirects)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal yaml")
	}

	redirectsMap := make(map[string]string)
	for _, entry := range redirects {
		redirectsMap[entry.Path] = entry.URL
	}
	return MapHandler(redirectsMap, fallback), nil
}
