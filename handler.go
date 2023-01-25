package handler

import (
	"context"
	"encoding/json"
	"github.com/graphql-go/graphql"
	"github.com/valyala/fasthttp"
	"net/http"
	"strings"
)

var (
	ContentTypeHTML           = "text/html"
	ContentTypeJSON           = "application/json"
	ContentTypeGraphQL        = "application/graphql"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
)

type Handler struct {
	Schema       *graphql.Schema
	pretty       bool
	graphiql     bool
	playground   bool
	rootObjectFn RootObjectFn
}

type RequestOptions struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

// a workaround for getting`variables` as a JSON string
type requestOptionsCompatibility struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

func getFromArgs(values *fasthttp.Args) *RequestOptions {
	query := values.Peek("query")
	if query != nil {
		// get variables map
		variables := make(map[string]interface{}, values.Len())
		variablesStr := values.Peek("variables")
		if variablesStr != nil {
			err := json.Unmarshal(variablesStr, &variables)
			if err != nil {
				return nil
			}
		}

		return &RequestOptions{
			Query:         string(query),
			Variables:     variables,
			OperationName: string(values.Peek("operationName")),
		}
	}

	return nil
}

// RequestOptions Parses a http.Request into GraphQL request options struct
func NewRequestOptions(ctx *fasthttp.RequestCtx) *RequestOptions {

	if reqOpt := getFromArgs(ctx.URI().QueryArgs()); reqOpt != nil {
		return reqOpt
	}

	if !ctx.Request.Header.IsPost() || len(ctx.Request.Body()) == 0 {
		return &RequestOptions{}
	}

	// TODO: improve Content-Type handling
	contentTypeStr := string(ctx.Request.Header.ContentType())
	contentTypeTokens := strings.Split(contentTypeStr, ";")
	contentType := contentTypeTokens[0]

	switch contentType {
	case ContentTypeGraphQL:
		body := ctx.Request.Body()
		return &RequestOptions{
			Query: string(body),
		}
	case ContentTypeFormURLEncoded:
		args := ctx.PostArgs()
		if args == nil {
			return &RequestOptions{}
		}

		if reqOpt := getFromArgs(args); reqOpt != nil {
			return reqOpt
		}

		return &RequestOptions{}

	case ContentTypeJSON:
		fallthrough
	default:
		var opts RequestOptions
		body := ctx.Request.Body()
		err := json.Unmarshal(body, &opts)
		if err != nil {
			// Probably `variables` was sent as a string instead of an object.
			// So, we try to be polite and try to parse that as a JSON string
			var optsCompatible requestOptionsCompatibility
			json.Unmarshal(body, &optsCompatible)
			json.Unmarshal([]byte(optsCompatible.Variables), &opts.Variables)
		}
		return &opts
	}
}

// ContextHandler provides an entrypoint into executing graphQL queries with a
// user-provided context.
func (h *Handler) ContextHandler(ctx context.Context, ctxreq *fasthttp.RequestCtx) {
	// get query
	opts := NewRequestOptions(ctxreq)

	// execute graphql query
	params := graphql.Params{
		Schema:         *h.Schema,
		RequestString:  opts.Query,
		VariableValues: opts.Variables,
		OperationName:  opts.OperationName,
		Context:        ctx,
	}
	if h.rootObjectFn != nil {
		params.RootObject = h.rootObjectFn(ctx, &ctxreq.Request)
	}
	result := graphql.Do(params)

	if h.graphiql {
		acceptHeader := string(ctxreq.Request.Header.Peek("Accept"))
		if !ctxreq.Request.URI().QueryArgs().Has("raw") && !strings.Contains(acceptHeader, ContentTypeJSON) && strings.Contains(acceptHeader, ContentTypeHTML) {
			renderGraphiQL(ctxreq, params)
			return
		}
	}

	if h.playground {
		acceptHeader := string(ctxreq.Request.Header.Peek("Accept"))
		if !ctxreq.Request.URI().QueryArgs().Has("raw") && !strings.Contains(acceptHeader, ContentTypeJSON) && strings.Contains(acceptHeader, ContentTypeHTML) {
			renderPlayground(ctxreq)
			return
		}
	}

	// use proper JSON Header
	ctxreq.Response.Header.SetContentType("application/json; charset=utf-8")

	if h.pretty {
		ctxreq.Response.SetStatusCode(http.StatusOK)
		buff, _ := json.MarshalIndent(result, "", "\t")

		ctxreq.Response.AppendBody(buff)
	} else {
		ctxreq.Response.SetStatusCode(http.StatusOK)
		buff, _ := json.Marshal(result)

		ctxreq.Response.AppendBody(buff)
	}
}

// ServeHTTP provides an entrypoint into executing graphQL queries.
func (h *Handler) ServeHTTP(ctx *fasthttp.RequestCtx) {
	h.ContextHandler(context.Background(), ctx)
}

// RootObjectFn allows a user to generate a RootObject per request
type RootObjectFn func(ctx context.Context, r *fasthttp.Request) map[string]interface{}

type Config struct {
	Schema       *graphql.Schema
	Pretty       bool
	GraphiQL     bool
	Playground   bool
	RootObjectFn RootObjectFn
}

func NewConfig() *Config {
	return &Config{
		Schema:     nil,
		Pretty:     true,
		GraphiQL:   true,
		Playground: false,
	}
}

func New(p *Config) *Handler {
	if p == nil {
		p = NewConfig()
	}
	if p.Schema == nil {
		panic("undefined GraphQL schema")
	}

	return &Handler{
		Schema:       p.Schema,
		pretty:       p.Pretty,
		graphiql:     p.GraphiQL,
		playground:   p.Playground,
		rootObjectFn: p.RootObjectFn,
	}
}
