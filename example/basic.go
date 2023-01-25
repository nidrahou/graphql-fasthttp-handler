package main

import (
	"github.com/graphql-go/graphql/testutil"
	handler "github.com/nidrahou/graphql-fasthttp-handler"
	"github.com/valyala/fasthttp"
)

func main() {
	h := handler.New(&handler.Config{
		Schema:   &testutil.StarWarsSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		h.ServeHTTP(ctx)
	})
}
