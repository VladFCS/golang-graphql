package graphql

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"

	"github.com/vladfc/ghira/internal/graphql/generated"
	"github.com/vladfc/ghira/internal/graphql/resolver"
)

func NewHandler(resolver *resolver.Resolver) http.Handler {
	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: resolver,
	})
	return handler.NewDefaultServer(schema)
}