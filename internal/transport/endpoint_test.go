package transport

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {

	type requestOptions struct {
		method string
		path   string
	}

	type testCase struct {
		name     string
		opt      requestOptions
		wantCode int
		wantBody string
	}

	ts := httptest.NewServer(app.router)
	defer ts.Close()

	testCases := []testCase{
		{
			name: "test 1",
			opt: requestOptions{
				method: http.MethodGet,
				path:   "/test-key",
			},
			wantCode: http.StatusOK,
			wantBody: "{\"data\":\"test-val\",\"result\":\"success\"}",
		},
		{
			name: "test 2",
			opt: requestOptions{
				method: http.MethodGet,
				path:   "/test-key-new",
			},
			wantCode: http.StatusOK,
			wantBody: "{\"data\":\"test-val-new\",\"result\":\"success\"}",
		},
		{
			name: "test bad",
			opt: requestOptions{
				method: http.MethodGet,
				path:   "/hsge",
			},
			wantCode: http.StatusNotFound,
			wantBody: "{\"data\":\"\",\"result\":\"not found\"}",
		},
		{
			name: "404 case",
			opt: requestOptions{
				method: http.MethodGet,
				path:   "/anything/qwerty",
			},
			wantCode: http.StatusNotFound,
			wantBody: "404 page not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.opt.method, tc.opt.path, nil)
			if err != nil {
				t.Log(err.Error())
			}
			w := httptest.NewRecorder()

			req.Header.Add("Content-Type", "application/json")
			app.router.ServeHTTP(w, req)

			assert.Equal(t, tc.wantCode, w.Code)

			assert.Equal(t, tc.wantBody, w.Body.String())
		})
	}

}
