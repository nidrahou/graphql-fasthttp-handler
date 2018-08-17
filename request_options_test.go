package handler_test

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/graphql-go/graphql/testutil"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	handler "."
)

func newCtx(method, u string, body io.Reader) (*fasthttp.RequestCtx, error) {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod(method)
	// parsedUrl, err := url.Parse(u)
	// if err != nil {
	// 	return nil, err
	// }
	// ctx.Request.Header.SetHost(parsedUrl.Host)
	ctx.Request.URI().Update(u)
	if body != nil {
		b, err := ioutil.ReadAll(body)
		if err != nil {
			return nil, err
		}
		ctx.Request.AppendBody(b)
	}
	return ctx, nil
}

func TestRequestOptions_GET_BasicQueryString(t *testing.T) {
	queryString := "query=query RebelsShipsQuery { rebels { name } }"
	expected := &handler.RequestOptions{
		Query:     "query RebelsShipsQuery { rebels { name } }",
		Variables: make(map[string]interface{}),
	}

	req, _ := newCtx("GET", fmt.Sprintf("/graphql?%v", queryString), nil)
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_GET_ContentTypeApplicationGraphQL(t *testing.T) {
	body := []byte(`query RebelsShipsQuery { rebels { name } }`)
	expected := &handler.RequestOptions{}

	req, _ := newCtx("GET", "/graphql", bytes.NewBuffer(body))
	req.Request.Header.SetContentType( "application/graphql")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_GET_ContentTypeApplicationJSON(t *testing.T) {
	body := `
	{
		"query": "query RebelsShipsQuery { rebels { name } }"
	}`
	expected := &handler.RequestOptions{}

	req, _ := newCtx("GET", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_GET_ContentTypeApplicationUrlEncoded(t *testing.T) {
	data := url.Values{}
	data.Add("query", "query RebelsShipsQuery { rebels { name } }")

	expected := &handler.RequestOptions{}

	req, _ := newCtx("GET", "/graphql", bytes.NewBufferString(data.Encode()))
	req.Request.Header.SetContentType( "application/x-www-form-urlencoded")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_BasicQueryString_WithNoBody(t *testing.T) {
	queryString := "query=query RebelsShipsQuery { rebels { name } }"
	expected := &handler.RequestOptions{
		Query:     "query RebelsShipsQuery { rebels { name } }",
		Variables: make(map[string]interface{}),
	}

	req, _ := newCtx("POST", fmt.Sprintf("/graphql?%v", queryString), nil)
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationGraphQL(t *testing.T) {
	body := []byte(`query RebelsShipsQuery { rebels { name } }`)
	expected := &handler.RequestOptions{
		Query: "query RebelsShipsQuery { rebels { name } }",
	}

	req, _ := newCtx("POST", "/graphql", bytes.NewBuffer(body))
	req.Request.Header.SetContentType("application/graphql")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationGraphQL_WithNonGraphQLQueryContent(t *testing.T) {
	body := []byte(`not a graphql query`)
	expected := &handler.RequestOptions{
		Query: "not a graphql query",
	}

	req, _ := newCtx("POST", "/graphql", bytes.NewBuffer(body))
	req.Request.Header.SetContentType( "application/graphql")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationGraphQL_EmptyBody(t *testing.T) {
	body := []byte(``)
	expected := &handler.RequestOptions{
		Query: "",
	}

	req, _ := newCtx("POST", "/graphql", bytes.NewBuffer(body))
	req.Request.Header.SetContentType( "application/graphql")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationGraphQL_NilBody(t *testing.T) {
	expected := &handler.RequestOptions{}

	req, _ := newCtx("POST", "/graphql", nil)
	req.Request.Header.SetContentType( "application/graphql")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationJSON(t *testing.T) {
	body := `
	{
		"query": "query RebelsShipsQuery { rebels { name } }"
	}`
	expected := &handler.RequestOptions{
		Query: "query RebelsShipsQuery { rebels { name } }",
	}

	req, _ := newCtx("POST", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_GET_WithVariablesAsObject(t *testing.T) {
	variables := url.QueryEscape(`{ "a": 1, "b": "2" }`)
	query := url.QueryEscape("query RebelsShipsQuery { rebels { name } }")
	queryString := fmt.Sprintf("query=%s&variables=%s", query, variables)
	expected := &handler.RequestOptions{
		Query: "query RebelsShipsQuery { rebels { name } }",
		Variables: map[string]interface{}{
			"a": float64(1),
			"b": "2",
		},
	}

	req, _ := newCtx("GET", fmt.Sprintf("/graphql?%v", queryString), nil)
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationJSON_WithVariablesAsObject(t *testing.T) {
	body := `
	{
		"query": "query RebelsShipsQuery { rebels { name } }",
		"variables": { "a": 1, "b": "2" }
	}`
	expected := &handler.RequestOptions{
		Query: "query RebelsShipsQuery { rebels { name } }",
		Variables: map[string]interface{}{
			"a": float64(1),
			"b": "2",
		},
	}

	req, _ := newCtx("POST", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationJSON_WithVariablesAsString(t *testing.T) {
	body := `
	{
		"query": "query RebelsShipsQuery { rebels { name } }",
		"variables": "{ \"a\": 1, \"b\": \"2\" }"
	}`
	expected := &handler.RequestOptions{
		Query: "query RebelsShipsQuery { rebels { name } }",
		Variables: map[string]interface{}{
			"a": float64(1),
			"b": "2",
		},
	}

	req, _ := newCtx("POST", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationJSON_WithInvalidJSON(t *testing.T) {
	body := `INVALIDJSON{}`
	expected := &handler.RequestOptions{}

	req, _ := newCtx("POST", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationJSON_WithNilBody(t *testing.T) {
	expected := &handler.RequestOptions{}

	req, _ := newCtx("POST", "/graphql", nil)
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationUrlEncoded(t *testing.T) {
	data := url.Values{}
	data.Add("query", "query RebelsShipsQuery { rebels { name } }")

	expected := &handler.RequestOptions{
		Query:     "query RebelsShipsQuery { rebels { name } }",
		Variables: make(map[string]interface{}),
	}

	req, _ := newCtx("POST", "/graphql", bytes.NewBufferString(data.Encode()))
	req.Request.Header.SetContentType( "application/x-www-form-urlencoded")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationUrlEncoded_WithInvalidData(t *testing.T) {
	data := "Invalid Data"

	expected := &handler.RequestOptions{}

	req, _ := newCtx("POST", "/graphql", bytes.NewBufferString(data))
	req.Request.Header.SetContentType( "application/x-www-form-urlencoded")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_ContentTypeApplicationUrlEncoded_WithNilBody(t *testing.T) {

	expected := &handler.RequestOptions{}

	req, _ := newCtx("POST", "/graphql", nil)
	req.Request.Header.SetContentType( "application/x-www-form-urlencoded")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_PUT_BasicQueryString(t *testing.T) {
	queryString := "query=query RebelsShipsQuery { rebels { name } }"
	expected := &handler.RequestOptions{
		Query:     "query RebelsShipsQuery { rebels { name } }",
		Variables: make(map[string]interface{}),
	}

	req, _ := newCtx("PUT", fmt.Sprintf("/graphql?%v", queryString), nil)
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_PUT_ContentTypeApplicationGraphQL(t *testing.T) {
	body := []byte(`query RebelsShipsQuery { rebels { name } }`)
	expected := &handler.RequestOptions{}

	req, _ := newCtx("PUT", "/graphql", bytes.NewBuffer(body))
	req.Request.Header.SetContentType( "application/graphql")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_PUT_ContentTypeApplicationJSON(t *testing.T) {
	body := `
	{
		"query": "query RebelsShipsQuery { rebels { name } }"
	}`
	expected := &handler.RequestOptions{}

	req, _ := newCtx("PUT", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_PUT_ContentTypeApplicationUrlEncoded(t *testing.T) {
	data := url.Values{}
	data.Add("query", "query RebelsShipsQuery { rebels { name } }")

	expected := &handler.RequestOptions{}

	req, _ := newCtx("PUT", "/graphql", bytes.NewBufferString(data.Encode()))
	req.Request.Header.SetContentType( "application/x-www-form-urlencoded")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_DELETE_BasicQueryString(t *testing.T) {
	queryString := "query=query RebelsShipsQuery { rebels { name } }"
	expected := &handler.RequestOptions{
		Query:     "query RebelsShipsQuery { rebels { name } }",
		Variables: make(map[string]interface{}),
	}

	req, _ := newCtx("DELETE", fmt.Sprintf("/graphql?%v", queryString), nil)
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_DELETE_ContentTypeApplicationGraphQL(t *testing.T) {
	body := []byte(`query RebelsShipsQuery { rebels { name } }`)
	expected := &handler.RequestOptions{}

	req, _ := newCtx("DELETE", "/graphql", bytes.NewBuffer(body))
	req.Request.Header.SetContentType( "application/graphql")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_DELETE_ContentTypeApplicationJSON(t *testing.T) {
	body := `
	{
		"query": "query RebelsShipsQuery { rebels { name } }"
	}`
	expected := &handler.RequestOptions{}

	req, _ := newCtx("DELETE", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/json")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_DELETE_ContentTypeApplicationUrlEncoded(t *testing.T) {
	data := url.Values{}
	data.Add("query", "query RebelsShipsQuery { rebels { name } }")

	expected := &handler.RequestOptions{}

	req, _ := newCtx("DELETE", "/graphql", bytes.NewBufferString(data.Encode()))
	req.Request.Header.SetContentType( "application/x-www-form-urlencoded")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestRequestOptions_POST_UnsupportedContentType(t *testing.T) {
	body := `<xml>query{}</xml>`
	expected := &handler.RequestOptions{}

	req, _ := newCtx("POST", "/graphql", bytes.NewBufferString(body))
	req.Request.Header.SetContentType( "application/xml")
	result := handler.NewRequestOptions(req)

	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
