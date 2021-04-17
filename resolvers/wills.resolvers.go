package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/kolioDev/after_life/graphql/model"
)

func (r *mutationResolver) CreateWill(ctx context.Context, willInput model.WillInput) (*model.Will, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Will(ctx context.Context) ([]*model.Will, error) {
	panic(fmt.Errorf("not implemented"))
}
