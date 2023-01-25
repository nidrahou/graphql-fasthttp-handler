package handler_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	handler "github.com/nidrahou/graphql-fasthttp-handler"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func decodeResponse(t *testing.T, ctx *fasthttp.RequestCtx) *graphql.Result {
	// clone request body reader so that we can have a nicer error message
	bodyString := ""
	var target graphql.Result
	bodyString = string(ctx.Response.Body())
	readerClone := strings.NewReader(bodyString)

	decoder := json.NewDecoder(readerClone)
	err := decoder.Decode(&target)
	if err != nil {
		t.Fatalf("DecodeResponseToType(): %v \n%v", err.Error(), bodyString)
	}
	return &target
}

func executeTest(t *testing.T, h *handler.Handler, ctx *fasthttp.RequestCtx) *graphql.Result {
	h.ServeHTTP(ctx)
	result := decodeResponse(t, ctx)
	return result
}

func newHTTPCtx(method, url string, body []byte) *fasthttp.RequestCtx {
	result := &fasthttp.RequestCtx{}
	result.Request.Header.SetMethod(method)
	result.Request.SetRequestURI(url)
	if body != nil {
		result.Request.AppendBody(body)
	}
	return result
}

func TestContextPropagated(t *testing.T) {
	myNameQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Name: "name",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return p.Context.Value("name"), nil
				},
			},
		},
	})
	myNameSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: myNameQuery,
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"name": "context-data",
		},
	}
	queryString := `query={name}`
	httpCtx := newHTTPCtx("GET", fmt.Sprintf("/graphql?%v", queryString), nil)

	h := handler.New(&handler.Config{
		Schema: &myNameSchema,
		Pretty: true,
	})

	ctx := context.WithValue(context.Background(), "name", "context-data")
	resp := httptest.NewRecorder()
	h.ContextHandler(ctx, httpCtx)
	result := decodeResponse(t, httpCtx)
	if resp.Code != http.StatusOK {
		t.Fatalf("unexpected server response %v", resp.Code)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestHandler_BasicQuery_Pretty(t *testing.T) {
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"hero": map[string]interface{}{
				"name": "R2-D2",
			},
		},
	}
	queryString := `query=query HeroNameQuery { hero { name } }`
	httpCtx := newHTTPCtx("GET", fmt.Sprintf("/graphql?%v", queryString), nil)

	h := handler.New(&handler.Config{
		Schema: &testutil.StarWarsSchema,
		Pretty: true,
	})
	result := executeTest(t, h, httpCtx)
	if httpCtx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("unexpected server response %v", httpCtx.Response.StatusCode())
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestHandler_BasicQuery_Ugly(t *testing.T) {
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"hero": map[string]interface{}{
				"name": "R2-D2",
			},
		},
	}
	queryString := `query=query HeroNameQuery { hero { name } }`
	httpCtx := newHTTPCtx("GET", fmt.Sprintf("/graphql?%v", queryString), nil)

	h := handler.New(&handler.Config{
		Schema: &testutil.StarWarsSchema,
		Pretty: false,
	})
	result := executeTest(t, h, httpCtx)
	if httpCtx.Response.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected server response %v", httpCtx.Response.StatusCode())
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestHandler_Params_NilParams(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if str, ok := r.(string); ok {
				if str != "undefined GraphQL schema" {
					t.Fatalf("unexpected error, got %v", r)
				}
				// test passed
				return
			}
			t.Fatalf("unexpected error, got %v", r)

		}
		t.Fatalf("expected to panic, did not panic")
	}()
	_ = handler.New(nil)

}

func TestHandler_BasicQuery_WithRootObjFn(t *testing.T) {
	myNameQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Name: "name",
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					rv := p.Info.RootValue.(map[string]interface{})
					return rv["rootValue"], nil
				},
			},
		},
	})
	myNameSchema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: myNameQuery,
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := &graphql.Result{
		Data: map[string]interface{}{
			"name": "foo",
		},
	}
	queryString := `query={name}`
	httpCtx := newHTTPCtx("GET", fmt.Sprintf("/graphql?%v", queryString), nil)

	h := handler.New(&handler.Config{
		Schema: &myNameSchema,
		Pretty: true,
		RootObjectFn: func(ctx context.Context, r *fasthttp.Request) map[string]interface{} {
			return map[string]interface{}{"rootValue": "foo"}
		},
	})
	result := executeTest(t, h, httpCtx)
	if httpCtx.Response.StatusCode() != http.StatusOK {
		t.Fatalf("unexpected server response %v", httpCtx.Response.StatusCode())
	}
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
