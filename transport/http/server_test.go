package http_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ThomasNguyenGitHub/go/endpoint"
	gotransport "github.com/ThomasNguyenGitHub/go/transport/http"
)

func TestServerBadDecode(t *testing.T) {
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, errors.New("dang") },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerBadEndpoint(t *testing.T) {
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, errors.New("dang") },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerBadEncode(t *testing.T) {
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, http.ResponseWriter, interface{}) error { return errors.New("dang") },
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusInternalServerError, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerErrorEncoder(t *testing.T) {
	errTeapot := errors.New("teapot")
	code := func(err error) int {
		if err == errTeapot {
			return http.StatusTeapot
		}
		return http.StatusInternalServerError
	}
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return struct{}{}, errTeapot },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(context.Context, http.ResponseWriter, interface{}) error { return nil },
		gotransport.ServerErrorEncoder(func(_ context.Context, err error, w http.ResponseWriter) { w.WriteHeader(code(err)) }),
	)
	server := httptest.NewServer(handler)
	defer server.Close()
	resp, _ := http.Get(server.URL)
	if want, have := http.StatusTeapot, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

func TestServerHappyPath(t *testing.T) {
	step, response := testServer(t)
	step()
	resp := <-response
	defer resp.Body.Close()
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := http.StatusOK, resp.StatusCode; want != have {
		t.Errorf("want %d, have %d (%s)", want, have, buf)
	}
}

func TestMultipleServerBefore(t *testing.T) {
	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := gotransport.NewServer(
		endpoint.Nop,
		func(context.Context, *http.Request) (interface{}, error) {
			return struct{}{}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		gotransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			ctx = context.WithValue(ctx, "one", 1)

			return ctx
		}),
		gotransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerBefores are used")
			}

			close(done)
			return ctx
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

func TestMultipleServerAfter(t *testing.T) {
	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := gotransport.NewServer(
		endpoint.Nop,
		func(context.Context, *http.Request) (interface{}, error) {
			return struct{}{}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		gotransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			ctx = context.WithValue(ctx, "one", 1)

			return ctx
		}),
		gotransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context {
			if _, ok := ctx.Value("one").(int); !ok {
				t.Error("Value was not set properly when multiple ServerAfters are used")
			}

			close(done)
			return ctx
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

func TestServerFinalizer(t *testing.T) {
	var (
		headerKey    = "X-Henlo-Lizer"
		headerVal    = "Helllo you stinky lizard"
		statusCode   = http.StatusTeapot
		responseBody = "go eat a fly ugly\n"
		done         = make(chan struct{})
	)
	handler := gotransport.NewServer(
		endpoint.Nop,
		func(context.Context, *http.Request) (interface{}, error) {
			return struct{}{}, nil
		},
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error {
			w.Header().Set(headerKey, headerVal)
			w.WriteHeader(statusCode)
			w.Write([]byte(responseBody))
			return nil
		},
		gotransport.ServerFinalizer(func(ctx context.Context, code int, _ *http.Request) {
			if want, have := statusCode, code; want != have {
				t.Errorf("StatusCode: want %d, have %d", want, have)
			}

			responseHeader := ctx.Value(gotransport.ContextKeyResponseHeaders).(http.Header)
			if want, have := headerVal, responseHeader.Get(headerKey); want != have {
				t.Errorf("%s: want %q, have %q", headerKey, want, have)
			}

			responseSize := ctx.Value(gotransport.ContextKeyResponseSize).(int64)
			if want, have := int64(len(responseBody)), responseSize; want != have {
				t.Errorf("response size: want %d, have %d", want, have)
			}

			close(done)
		}),
	)

	server := httptest.NewServer(handler)
	defer server.Close()
	go http.Get(server.URL)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for finalizer")
	}
}

type enhancedResponse struct {
	Foo string `json:"foo"`
}

func (e enhancedResponse) StatusCode() int      { return http.StatusPaymentRequired }
func (e enhancedResponse) Headers() http.Header { return http.Header{"X-Edward": []string{"Snowden"}} }

func TestEncodeJSONResponse(t *testing.T) {
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return enhancedResponse{Foo: "bar"}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gotransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := http.StatusPaymentRequired, resp.StatusCode; want != have {
		t.Errorf("StatusCode: want %d, have %d", want, have)
	}
	if want, have := "Snowden", resp.Header.Get("X-Edward"); want != have {
		t.Errorf("X-Edward: want %q, have %q", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := `{"foo":"bar"}`, strings.TrimSpace(string(buf)); want != have {
		t.Errorf("Body: want %s, have %s", want, have)
	}
}

type multiHeaderResponse struct{}

func (_ multiHeaderResponse) Headers() http.Header {
	return http.Header{"Vary": []string{"Origin", "User-Agent"}}
}

func TestAddMultipleHeaders(t *testing.T) {
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return multiHeaderResponse{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gotransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	expect := map[string]map[string]struct{}{"Vary": map[string]struct{}{"Origin": struct{}{}, "User-Agent": struct{}{}}}
	for k, vls := range resp.Header {
		for _, v := range vls {
			delete((expect[k]), v)
		}
		if len(expect[k]) != 0 {
			t.Errorf("Header: unexpected header %s: %v", k, expect[k])
		}
	}
}

type multiHeaderResponseError struct {
	multiHeaderResponse
	msg string
}

func (m multiHeaderResponseError) Error() string {
	return m.msg
}

func TestAddMultipleHeadersErrorEncoder(t *testing.T) {
	errStr := "oh no"
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) {
			return nil, multiHeaderResponseError{msg: errStr}
		},
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gotransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	expect := map[string]map[string]struct{}{"Vary": map[string]struct{}{"Origin": struct{}{}, "User-Agent": struct{}{}}}
	for k, vls := range resp.Header {
		for _, v := range vls {
			delete((expect[k]), v)
		}
		if len(expect[k]) != 0 {
			t.Errorf("Header: unexpected header %s: %v", k, expect[k])
		}
	}
	if b, _ := ioutil.ReadAll(resp.Body); errStr != string(b) {
		t.Errorf("ErrorEncoder: got: %q, expected: %q", b, errStr)
	}
}

type noContentResponse struct{}

func (e noContentResponse) StatusCode() int { return http.StatusNoContent }

func TestEncodeNoContent(t *testing.T) {
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return noContentResponse{}, nil },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		gotransport.EncodeJSONResponse,
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if want, have := http.StatusNoContent, resp.StatusCode; want != have {
		t.Errorf("StatusCode: want %d, have %d", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := 0, len(buf); want != have {
		t.Errorf("Body: want no content, have %d bytes", have)
	}
}

type enhancedError struct{}

func (e enhancedError) Error() string                { return "enhanced error" }
func (e enhancedError) StatusCode() int              { return http.StatusTeapot }
func (e enhancedError) MarshalJSON() ([]byte, error) { return []byte(`{"err":"enhanced"}`), nil }
func (e enhancedError) Headers() http.Header         { return http.Header{"X-Enhanced": []string{"1"}} }

func TestEnhancedError(t *testing.T) {
	handler := gotransport.NewServer(
		func(context.Context, interface{}) (interface{}, error) { return nil, enhancedError{} },
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
		func(_ context.Context, w http.ResponseWriter, _ interface{}) error { return nil },
	)

	server := httptest.NewServer(handler)
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if want, have := http.StatusTeapot, resp.StatusCode; want != have {
		t.Errorf("StatusCode: want %d, have %d", want, have)
	}
	if want, have := "1", resp.Header.Get("X-Enhanced"); want != have {
		t.Errorf("X-Enhanced: want %q, have %q", want, have)
	}
	buf, _ := ioutil.ReadAll(resp.Body)
	if want, have := `{"err":"enhanced"}`, strings.TrimSpace(string(buf)); want != have {
		t.Errorf("Body: want %s, have %s", want, have)
	}
}

func TestNoOpRequestDecoder(t *testing.T) {
	resw := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error("Failed to create request")
	}
	handler := gotransport.NewServer(
		func(ctx context.Context, request interface{}) (interface{}, error) {
			if request != nil {
				t.Error("Expected nil request in endpoint when using NopRequestDecoder")
			}
			return nil, nil
		},
		gotransport.NopRequestDecoder,
		gotransport.EncodeJSONResponse,
	)
	handler.ServeHTTP(resw, req)
	if resw.Code != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resw.Code)
	}
}

func testServer(t *testing.T) (step func(), resp <-chan *http.Response) {
	var (
		stepch   = make(chan bool)
		endpoint = func(context.Context, interface{}) (interface{}, error) { <-stepch; return struct{}{}, nil }
		response = make(chan *http.Response)
		handler  = gotransport.NewServer(
			endpoint,
			func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil },
			func(context.Context, http.ResponseWriter, interface{}) error { return nil },
			gotransport.ServerBefore(func(ctx context.Context, r *http.Request) context.Context { return ctx }),
			gotransport.ServerAfter(func(ctx context.Context, w http.ResponseWriter) context.Context { return ctx }),
		)
	)
	go func() {
		server := httptest.NewServer(handler)
		defer server.Close()
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Error(err)
			return
		}
		response <- resp
	}()
	return func() { stepch <- true }, response
}
