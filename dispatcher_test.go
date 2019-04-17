package covercheck_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/taiyoh/covercheck"
)

func TestDispatcher(t *testing.T) {
	mux := http.NewServeMux()
	testserver := httptest.NewServer(mux)
	defer testserver.Close()
	parts := strings.SplitN(testserver.URL, "/", 3)
	d := covercheck.NewDispatcher(parts[2])

	i := uint32(0)
	mux.HandleFunc("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("!!!!!"))
	})
	mux.HandleFunc("/hoge/fuga", func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&i, 1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hogefuga"))
	})

	d.AddController(covercheck.NewController("/healthcheck", func(context.Context) error {
		cli := &http.Client{}
		resp, err := cli.Get(testserver.URL + "/hoge/fuga")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		return nil
	}))

	proxy := httptest.NewServer(d.Mux())
	defer proxy.Close()

	cli := proxy.Client()

	{
		resp, err := cli.Get(proxy.URL + "/foo/bar")
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		b := bytes.NewBuffer([]byte{})
		b.ReadFrom(resp.Body)
		if body := b.String(); body != "!!!!!" {
			t.Errorf("GET /foo/bar response body is wrong: %s", body)
		}
		if i != uint32(0) {
			t.Error("arleady request comes.")
		}
	}
	{
		resp, err := cli.Get(proxy.URL + "/healthcheck")
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		b := bytes.NewBuffer([]byte{})
		b.ReadFrom(resp.Body)
		if body := b.String(); body != "ok" {
			t.Errorf("GET /foo/bar response body is wrong: %s", body)
		}
		if i != uint32(1) {
			t.Error("request not yet comes...")
		}
	}
}
