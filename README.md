[![CircleCI](https://circleci.com/gh/lab259/graphql-fasthttp-handler.svg?style=shield)](https://circleci.com/gh/lab259/graphql-fasthttp-handler)
[![codecov](https://codecov.io/gh/lab259/graphql-fasthttp-handler/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/graphql-fasthttp-handler)
[![GoDoc](https://godoc.org/lab259/graphql-fasthttp-handler?status.svg)](https://godoc.org/github.com/lab259/graphql-fasthttp-handler)

# graphql-go-handler

Golang HTTP.Handler for [graphl-go](https://github.com/graphql-go/graphql).

The original implementation by [graphql-go/handler](https://github.com/graphql-go/handler)
uses the default HTTP implementation. This fork adapts the library for using the
[valyala/fasthttp](https://github.com/valyala/fasthttp).

### Usage

```go
package main

import (
	"github.com/valyala/fasthttp"
	handler "github.com/lab259/graphql-fashttp-handler"
	"github.com/graphql-go/graphql/testutil"
)

func main() {
	h := handler.New(&handler.Config{
		Schema:   &testutil.StarWarsSchema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Serving given endpoint for the sake of example
	fasthttp.ListenAndServe(":8080", func(ctx *fasthttp.RequestCtx) {
		h.ServeHTTP(ctx)
	})
}

```


### Details

The handler will accept requests with the parameters:

  * **`query`**: A string GraphQL document to be executed.

  * **`variables`**: The runtime values to use for any GraphQL query variables
    as a JSON object.

  * **`operationName`**: If the provided `query` contains multiple named
    operations, this specifies which operation should be executed. If not
    provided, an 400 error will be returned if the `query` contains multiple
    named operations.

GraphQL will first look for each parameter in the URL's query-string:

```
/graphql?query=query+getUser($id:ID){user(id:$id){name}}&variables={"id":"4"}
```

If not found in the query-string, it will look in the POST request body.
The `handler` will interpret it
depending on the provided `Content-Type` header.

  * **`application/json`**: the POST body will be parsed as a JSON
    object of parameters.

  * **`application/x-www-form-urlencoded`**: this POST body will be
    parsed as a url-encoded string of key-value pairs.

  * **`application/graphql`**: The POST body will be parsed as GraphQL
    query string, which provides the `query` parameter.


### Examples
- [golang-graphql-playground](https://github.com/graphql-go/playground)
- [golang-relay-starter-kit](https://github.com/sogko/golang-relay-starter-kit)
- [todomvc-relay-go](https://github.com/sogko/todomvc-relay-go)

### Test

```bash
$ go get -u github.com/onsi/ginkgo/ginkgo
$ make deps
$ make test
```

### Credits

This project is originally forked from [graphql-go/handler](https://github.com/graphql-go/handler).

### License

MIT