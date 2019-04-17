package covercheck

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// RequestChecker provides checker interface using HTTP request.
type RequestChecker struct {
	Method  string
	URL     *url.URL
	Headers http.Header
	Body    io.Reader
}

// NewRequestCheckerGET returns RequestChecker object.
func NewRequestCheckerGET(u *url.URL, headers http.Header) *RequestChecker {
	if headers == nil {
		headers = map[string][]string{}
	}
	return &RequestChecker{http.MethodGet, u, headers, nil}
}

// Checker returns Checker typed operation.
func (c *RequestChecker) Checker() Checker {
	return func(ctx context.Context) error {
		req, err := http.NewRequest(c.Method, c.URL.String(), c.Body)
		if err != nil {
			return err
		}
		for name, values := range c.Headers {
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
