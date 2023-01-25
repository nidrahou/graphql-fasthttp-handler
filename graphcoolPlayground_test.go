package handler_test

import (
	"fmt"
	"github.com/graphql-go/graphql/testutil"
	handler "github.com/nidrahou/graphql-fasthttp-handler"
	"net/http"
	"strings"
	"testing"
)

func TestRenderPlayground(t *testing.T) {
	cases := map[string]struct {
		playgroundEnabled    bool
		accept               string
		url                  string
		expectedStatusCode   int
		expectedContentType  string
		expectedBodyContains string
	}{
		"renders Playground": {
			playgroundEnabled:    true,
			accept:               "text/html",
			expectedStatusCode:   http.StatusOK,
			expectedContentType:  "text/html; charset=utf-8",
			expectedBodyContains: "<!DOCTYPE html>",
		},
		"doesn't render Playground if turned off": {
			playgroundEnabled:   false,
			accept:              "text/html",
			expectedStatusCode:  http.StatusOK,
			expectedContentType: "application/json; charset=utf-8",
		},
		"doesn't render Playground if Content-Type application/json is present": {
			playgroundEnabled:   true,
			accept:              "application/json,text/html",
			expectedStatusCode:  http.StatusOK,
			expectedContentType: "application/json; charset=utf-8",
		},
		"doesn't render Playground if Content-Type text/html is not present": {
			playgroundEnabled:   true,
			expectedStatusCode:  http.StatusOK,
			expectedContentType: "application/json; charset=utf-8",
		},
		"doesn't render Playground if 'raw' query is present": {
			playgroundEnabled:   true,
			accept:              "text/html",
			url:                 "?raw",
			expectedStatusCode:  http.StatusOK,
			expectedContentType: "application/json; charset=utf-8",
		},
	}

	for tcID, tc := range cases {
		t.Run(tcID, func(t *testing.T) {
			ctx := newHTTPCtx("GET", tc.url, nil)
			fmt.Println(tc.url)
			ctx.Request.Header.Set("Accept", tc.accept)

			h := handler.New(&handler.Config{
				Schema:     &testutil.StarWarsSchema,
				GraphiQL:   false,
				Playground: tc.playgroundEnabled,
			})

			h.ServeHTTP(ctx)

			statusCode := ctx.Response.StatusCode()
			if statusCode != tc.expectedStatusCode {
				t.Fatalf("%s: wrong status code, expected %v, got %v", tcID, tc.expectedStatusCode, statusCode)
			}

			contentType := ctx.Response.Header.Peek("Content-Type")
			if string(contentType) != tc.expectedContentType {
				t.Fatalf("%s: wrong content type, expected %s, got %s", tcID, tc.expectedContentType, contentType)
			}

			body := string(ctx.Response.Body())
			if !strings.Contains(body, tc.expectedBodyContains) {
				t.Fatalf("%s: wrong body, expected %s to contain %s", tcID, body, tc.expectedBodyContains)
			}
		})
	}
}
