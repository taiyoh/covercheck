package covercheck_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/taiyoh/covercheck"
)

func TestRequestChceckGET(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-For") != "forwarder" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("NG"))
			return
		}
		if r.URL.RawQuery != "hoge=fuga" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("NG"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	testserver := httptest.NewServer(mux)
	defer testserver.Close()

	for _, tt := range []struct {
		label        string
		forwardedFor string
		query        string
		expectedErr  bool
	}{
		{
			"wrong query",
			"forwarder",
			"hoge=fugo",
			true,
		},
		{
			"wrong forwardfor",
			"",
			"hoge=fuga",
			true,
		},
		{
			"valid request",
			"forwarder",
			"hoge=fuga",
			false,
		},
	} {
		u, _ := url.Parse(testserver.URL + "/foo?" + tt.query)
		c := covercheck.NewRequestCheckerGET(tt.forwardedFor, u, nil)
		checker := c.Checker()
		actual := checker(context.Background()) != nil
		if actual != tt.expectedErr {
			t.Errorf("[%s] expected:%v, actual:%v", tt.label, tt.expectedErr, actual)
		}
	}
}

func TestRequestChceckPOST(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.RawQuery != "hoge=fuga" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("NG"))
			return
		}
		r.ParseForm()
		vals := r.PostForm
		if vals.Get("foo") != "bar" || vals.Get("aaa") != "iii" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("NG"))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	testserver := httptest.NewServer(mux)
	defer testserver.Close()

	headers := http.Header{
		"Content-Type": []string{
			"application/x-www-form-urlencoded",
		},
	}
	for _, tt := range []struct {
		label       string
		query       string
		form        string
		expectedErr bool
	}{
		{
			"wrong query",
			"hoge=fugo",
			"foo=bar&aaa=iii",
			true,
		},
		{
			"wrong post body",
			"hoge=fuga",
			"aaa=iii",
			true,
		},
		{
			"valid request",
			"hoge=fuga",
			"foo=bar&aaa=iii",
			false,
		},
	} {
		u, _ := url.Parse(testserver.URL + "/foo?" + tt.query)
		c := covercheck.NewRequestCheckerPOST("forwarder-test", u, headers, []byte(tt.form))
		checker := c.Checker()
		{
			actual := checker(context.Background()) != nil
			if actual != tt.expectedErr {
				t.Errorf("[%s] expected:%v, actual:%v", tt.label, tt.expectedErr, actual)
			}
		}
		{
			actual := checker(context.Background()) != nil
			if actual != tt.expectedErr {
				t.Errorf("2nd [%s] expected:%v, actual:%v", tt.label, tt.expectedErr, actual)
			}
		}
	}
}
