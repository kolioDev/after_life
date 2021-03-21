package actions

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/gobuffalo/buffalo"
	"github.com/kolioDev/after_life/graphql/generated"
	"github.com/kolioDev/after_life/models"
	"github.com/kolioDev/after_life/resolvers"
)

// GraphqlIndex default implementation.
func GraphqlIndex(c buffalo.Context) error {
	u := c.Value("current_user").(*models.User)
	tx := models.DB //c.Value("tx").(*pop.Connection)

	resolvers.SetUser(u)
	resolvers.SetTX(tx)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &resolvers.Resolver{},
	}))
	srv.ServeHTTP(c.Response(), c.Request())

	return nil
}

func GraphqlPlayground(c buffalo.Context) error {
	handler := playground.Handler("GraphQL playground", "/query")
	handler(c.Response(), c.Request())
	return nil
}
