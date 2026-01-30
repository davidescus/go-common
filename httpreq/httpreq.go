package httpreq

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	stdUrl "net/url"
)

// New creates a new http.Request with the specified method, URL, query parameters, and form parameters.
func New(ctx context.Context, method, url string, queryParams, formParams map[string]any) (*http.Request, error) {
	query := stdUrl.Values{}
	for k, v := range queryParams {
		query.Add(k, fmt.Sprintf("%v", v))
	}
	queryString := query.Encode()

	form := stdUrl.Values{}
	for k, v := range formParams {
		form.Set(k, fmt.Sprintf("%v", v))
	}
	bodyString := form.Encode()

	header := http.Header{}
	body := &bytes.Buffer{}
	if bodyString != "" {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		body = bytes.NewBufferString(bodyString)
	}

	if queryString != "" {
		url = fmt.Sprintf("%s?%s", url, queryString)
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return &http.Request{}, err
	}

	request = request.WithContext(ctx)
	request.Header = header

	return request, nil
}
