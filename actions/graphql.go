package actions

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gobuffalo/buffalo"
	"github.com/kolioDev/after_life/graphql/graph"
	"github.com/kolioDev/after_life/graphql/graph/generated"
	"github.com/kolioDev/after_life/models"
)

// GraphqlIndex default implementation.
func GraphqlIndex(c buffalo.Context) error {
	u := c.Value("current_user").(*models.User)

	graph.SetUser(u)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers: &graph.Resolver{},
	}))
	srv.ServeHTTP(c.Response(), c.Request())

	return nil
}
