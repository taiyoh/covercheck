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

func fillXFFHeader(f string, h http.Header) http.Header {
	if h == nil {
		h = http.Header{}
	}
	if xff, exists := h[xffKey]; exists {
		for i, v := range xff {
			xff[i] = fmt.Sprintf("%s, %s", v, f)
		}
		h[xffKey] = xff
		return h
	}
	h.Set(xffKey, f)
	return h
}

// NewRequestCheckerGET returns RequestChecker object.
func NewRequestCheckerGET(forwardedFor string, u *url.URL, headers http.Header) *RequestChecker {
	headers = fillXFFHeader(forwardedFor, headers)
	return &RequestChecker{http.MethodGet, u, headers, nil}
}

// NewRequestCheckerPOST returns RequestChecker object.
func NewRequestCheckerPOST(forwardedFor string, u *url.URL, headers http.Header, body []byte) *RequestChecker {
	headers = fillXFFHeader(forwardedFor, headers)
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
