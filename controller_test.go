package covercheck_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/taiyoh/covercheck"
)

func TestController(t *testing.T) {
	okChecker := func(context.Context) error {
		time.Sleep(time.Millisecond)
		return nil
	}
	ngChecker := func(context.Context) error {
		time.Sleep(time.Millisecond)
		return errors.New("ng")
	}

	ctrl := covercheck.NewController("/foo/bar/baz", okChecker)

	mux := http.NewServeMux()
	mux.HandleFunc(ctrl.HandlerFunc())

	testserver := httptest.NewServer(mux)
	defer testserver.Close()

	cli := testserver.Client()
	{
		resp, err := cli.Get(testserver.URL + "/foo/bar/baz")
		if err != nil {
			t.Errorf("GET /foo/bar/baz error: %v", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("GET /foo/bar/baz response fail: %d", resp.StatusCode)
		}
	}
	ctrl.AddChecker(ngChecker)
	{
		resp, err := cli.Get(testserver.URL + "/foo/bar/baz")
		if err != nil {
			t.Errorf("GET /foo/bar/baz error: %v", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("GET /foo/bar/baz response fail: %d", resp.StatusCode)
		}
	}
}
