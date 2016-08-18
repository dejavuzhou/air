package air

import (
	"bytes"
	"io"
	"mime/multipart"

	"github.com/valyala/fasthttp"
)

// Request represents the current HTTP request.
type Request struct {
	fastCtx *fasthttp.RequestCtx
	air     *Air

	Header *RequestHeader
	URI    *URI
}

// newRequest returns a new instance of `Request`.
func newRequest(a *Air) *Request {
	return &Request{
		air: a,
	}
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (r *Request) IsTLS() bool {
	return r.fastCtx.IsTLS()
}

// Scheme returns the HTTP protocol scheme, "http" or "https".
func (r *Request) Scheme() string {
	return string(r.fastCtx.Request.URI().Scheme())
}

// Host returns HTTP request host. Per RFC 2616, this is either the value of
// the "Host" header or the host name given in the URI itself.
func (r *Request) Host() string {
	return string(r.fastCtx.Request.Host())
}

// Referer returns the referring URI, if sent in the request.
func (r *Request) Referer() string {
	return string(r.fastCtx.Request.Header.Referer())
}

// ContentLength returns the size of request's body.
func (r *Request) ContentLength() int64 {
	return int64(r.fastCtx.Request.Header.ContentLength())
}

// UserAgent returns the client's "User-Agent".
func (r *Request) UserAgent() string {
	return string(r.fastCtx.UserAgent())
}

// RemoteAddress returns the client's network address.
func (r *Request) RemoteAddress() string {
	return r.fastCtx.RemoteAddr().String()
}

// RemoteIP returns the client's network ip address.
func (r *Request) RemoteIP() string {
	return r.fastCtx.RemoteIP().String()
}

// Method returns the request's HTTP function.
func (r *Request) Method() string {
	return string(r.fastCtx.Method())
}

// SetMethod sets the HTTP method of the request.
func (r *Request) SetMethod(method string) {
	r.fastCtx.Request.Header.SetMethodBytes([]byte(method))
}

// RequestURI returns the unmodified "Request-URI" sent by the client.
func (r *Request) RequestURI() string {
	return string(r.fastCtx.Request.RequestURI())
}

// SetRequestURI sets the "Request-URI".
func (r *Request) SetRequestURI(uri string) {
	r.fastCtx.Request.Header.SetRequestURI(uri)
}

// Body returns request's body.
func (r *Request) Body() io.Reader {
	return bytes.NewBuffer(r.fastCtx.Request.Body())
}

// SetBody sets request's body.
func (r *Request) SetBody(reader io.Reader) {
	r.fastCtx.Request.SetBodyStream(reader, 0)
}

// FormValue returns the form field value for the provided name.
func (r *Request) FormValue(name string) string {
	return string(r.fastCtx.FormValue(name))
}

// FormParams returns the form parameters.
func (r *Request) FormParams() (params map[string][]string) {
	params = make(map[string][]string)
	mf, err := r.fastCtx.Request.MultipartForm()

	if err == fasthttp.ErrNoMultipartForm {
		r.fastCtx.PostArgs().VisitAll(func(k, v []byte) {
			key := string(k)
			if _, ok := params[key]; ok {
				params[key] = append(params[key], string(v))
			} else {
				params[string(k)] = []string{string(v)}
			}
		})
	} else if err == nil {
		for k, v := range mf.Value {
			if len(v) > 0 {
				params[k] = v
			}
		}
	}

	return
}

// FormFile returns the multipart form file for the provided name.
func (r *Request) FormFile(name string) (*multipart.FileHeader, error) {
	return r.fastCtx.FormFile(name)
}

// MultipartForm returns the multipart form.
func (r *Request) MultipartForm() (*multipart.Form, error) {
	return r.fastCtx.MultipartForm()
}

// Cookie returns the named cookie provided in the request.
func (r *Request) Cookie(name string) (Cookie, error) {
	c := &fasthttp.Cookie{}
	b := r.fastCtx.Request.Header.Cookie(name)
	if b == nil {
		return Cookie{}, ErrCookieNotFound
	}
	c.SetKey(name)
	c.SetValueBytes(b)
	return Cookie{c}, nil
}

// Cookies returns the HTTP cookies sent with the request.
func (r *Request) Cookies() []Cookie {
	cookies := []Cookie{}
	r.fastCtx.Request.Header.VisitAllCookie(func(name, value []byte) {
		c := &fasthttp.Cookie{}
		c.SetKeyBytes(name)
		c.SetValueBytes(value)
		cookies = append(cookies, Cookie{c})
	})
	return cookies
}

// reset resets the instacne of `Request`.
func (r *Request) reset() {
	r.fastCtx = nil
	r.Header = nil
	r.URI = nil
}
