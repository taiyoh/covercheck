package covercheck

import (
	"io"
	"net/http"
)

// Dispatcher provides URL routing whether request is healthcheck or not.
type Dispatcher struct {
	backend string
	mux     *http.ServeMux
}

func cloneRequest(host string, r *http.Request) (*http.Request, error) {
	req, err := http.NewRequest(r.Method, "", r.Body)
	if err != nil {
		return nil, err
	}
	for name, values := range r.Header {
		for _, val := range values {
			req.Header.Set(name, val)
		}
	}
	u := r.URL
	u.Scheme = "http"
	u.Host = host
	req.URL = u
	req.RemoteAddr = r.RemoteAddr
	return req.WithContext(req.Context()), nil
}

func copyResponse(w http.ResponseWriter, resp *http.Response) {
	for name, values := range resp.Header {
		for _, val := range values {
			w.Header().Set(name, val)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// NewDispatcher returns Dispatcher object.
func NewDispatcher(backend string) *Dispatcher {
	mux := http.NewServeMux()
	d := &Dispatcher{
		backend: backend,
		mux:     mux,
	}
	d.mux.HandleFunc("/", d.ProxyHandleFunc)
	return d
}

// ProxyHandleFunc provides normal request handling for proxy.
func (d *Dispatcher) ProxyHandleFunc(w http.ResponseWriter, r *http.Request) {
	req, err := cloneRequest(d.backend, r)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	// combined for GET/POST
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	copyResponse(w, resp)
}

// AddController provides filling interface for add healthcheck path and checklist.
func (d *Dispatcher) AddController(c *Controller) {
	d.mux.HandleFunc(c.HandlerFunc())
}

// PathCheck provides syntax sugar for cc.AddController(covercheck.NewController("/foopath", checker1, checker2)).
func (d *Dispatcher) PathCheck(path string, checklist ...Checker) {
	d.AddController(NewController(path, checklist...))
}

// Mux returns http.ServeMux object
func (d *Dispatcher) Mux() *http.ServeMux {
	return d.mux
}
