package covercheck_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/taiyoh/covercheck"
)

func TestRequestChceckGET(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("OK: %s", r.URL.RawQuery)))
	})
	mux.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("NG"))
	})

	testserver := httptest.NewServer(mux)
	defer testserver.Close()

	{
		u, _ := url.Parse(testserver.URL + "/foo?hoge=fuga")
		c := covercheck.NewRequestCheckerGET(u, nil)
		checker := c.Checker()
		if err := checker(context.Background()); err != nil {
			t.Error(err)
		}
	}
	{
		u, _ := url.Parse(testserver.URL + "/bar?hoge=fuga")
		c := covercheck.NewRequestCheckerGET(u, nil)
		checker := c.Checker()
		if err := checker(context.Background()); err == nil {
			t.Error("status is wrong. error should be returned")
		}
	}
}
