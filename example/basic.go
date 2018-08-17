package main

import (
	"github.com/valyala/fasthttp"
	".."
	"github.com/graphql-go/graphql/testutil"
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
