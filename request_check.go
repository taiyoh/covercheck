package covercheck

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// RequestChecker provides checker interface using HTTP request.
type RequestChecker struct {
	method  string
	url     *url.URL
	headers http.Header
	body    []byte
}

const xffKey = "X-Forwarded-For"

// NewRequestCheckerGET returns RequestChecker object.
func NewRequestCheckerGET(forwardedFor string, u *url.URL, headers http.Header) *RequestChecker {
	if headers == nil {
		headers = http.Header{}
	}
	headers.Set(xffKey, forwardedFor)
	return &RequestChecker{http.MethodGet, u, headers, nil}
}

// NewRequestCheckerPOST returns RequestChecker object.
func NewRequestCheckerPOST(forwardedFor string, u *url.URL, headers http.Header, body []byte) *RequestChecker {
	if headers == nil {
		headers = http.Header{}
	}
	headers.Set(xffKey, forwardedFor)
	return &RequestChecker{http.MethodPost, u, headers, body}
}

// Checker returns Checker typed operation.
func (c *RequestChecker) Checker() Checker {
	return func(ctx context.Context) error {
		req, err := http.NewRequest(c.method, c.url.String(), bytes.NewBuffer(c.body))
		if err != nil {
			return err
		}
		for name, values := range c.headers {
			for _, val := range values {
				req.Header.Set(name, val)
			}
		}
		cli := &http.Client{}
		resp, err := cli.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("return status: %d", resp.StatusCode)
		}
		return nil
	}
}
