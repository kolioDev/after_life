package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gobuffalo/nulls"
	"github.com/kolioDev/after_life/graphql/model"
	"github.com/kolioDev/after_life/models"
	"github.com/kolioDev/after_life/scalars"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateWill(ctx context.Context, willInput model.WillInput) (*model.Will, error) {
	var instructions models.Instructions

	for _, i := range willInput.Instructions {
		instructions = append(instructions, models.Instruction{
			Text:  i.Text,
			Index: uint(i.Index),
		})
	}

	w := models.Will{
		Title: willInput.Title,
		Priority: nulls.UInt32{
			Valid:  GetNullableInt(willInput.Priority).Valid,
			UInt32: uint32(GetNullableInt(willInput.Priority).Int),
		},
		Instructions: &instructions,
	}

	verrs, err := w.Create(TX, User)
	if err != nil {
		return nil, err
	}
	if verrs.HasAny() {

		var extensions = map[string]interface{}{}

		for key, errs := range verrs.Errors {
			extensions[key] = errs[0]
		}

		graphql.AddError(ctx, &gqlerror.Error{
			Message:    "Validation errors",
			Extensions: extensions,
		})
		return nil, nil
	}

	return w.ToGraphQL(), nil
}

func (r *queryResolver) Will(ctx context.Context, id scalars.UUID) (*model.Will, error) {
	will := &models.Will{}

	err := will.Get(TX, scalars.GhqlUUID2ModelsUUID(id))

	if err != nil {
		return nil, err
	}

	if will.UserID != User.ID {
		graphql.AddError(ctx, &gqlerror.Error{
			Message: "User does not own will",
		})
		return nil, nil
	}

	return will.ToGraphQL(), nil
}
