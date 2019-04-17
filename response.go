package covercheck

import "net/http"

// Response provides healthcheck response.
type Response struct {
	code        int
	contentType string
	renderer    func() []byte
}

// Render provides filling response by type.
func (r Response) Render(w http.ResponseWriter) {
	w.Header().Set("Content-Type", r.contentType)
	w.WriteHeader(r.code)
	w.Write(r.renderer())
}

func newSuccess() Response {
	return Response{
		code:        http.StatusOK,
		contentType: "text/plain",
		renderer: func() []byte {
			return []byte("ok")
		},
	}
}

func newFail() Response {
	return Response{
		code:        http.StatusOK,
		contentType: "text/plain",
		renderer: func() []byte {
			return []byte("ng")
		},
	}
}
