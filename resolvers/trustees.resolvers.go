package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"database/sql"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gobuffalo/nulls"
	"github.com/kolioDev/after_life/graphql/model"
	"github.com/kolioDev/after_life/models"
	"github.com/kolioDev/after_life/scalars"
	errs "github.com/pkg/errors"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func (r *mutationResolver) CreateTrustee(ctx context.Context, trusteeInput model.TrusteeInput) (*model.Trustee, error) {
	t := models.Trustee{
		Name:           trusteeInput.Name,
		Relationship:   trusteeInput.Relationship.String(),
		Phone:          trusteeInput.Phone,
		Email:          trusteeInput.Email,
		FacebookLink:   GetNullable(trusteeInput.FacebookLink),
		TwitterLink:    GetNullable(trusteeInput.TwitterLink),
		AdditionalInfo: GetNullable(trusteeInput.AdditionalInformation),
	}

	verrs, err := t.Create(TX, User.ID)
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

	return t.ToGraphQL(), nil
}

func (r *mutationResolver) UpdateTrustee(ctx context.Context, trusteeInput model.UpdateTrustee) (*model.Trustee, error) {
	t := &models.Trustee{}
	if err := TX.Find(t, trusteeInput.ID); err != nil {
		if errs.Cause(err) != sql.ErrNoRows {
			graphql.AddError(ctx, &gqlerror.Error{
				Message: "Cannot find trustee",
			})
			return nil, nil
		}
		return nil, errs.WithStack(err)
	}

	if trusteeInput.Name != nil {
		t.Name = *trusteeInput.Name
	}

	if trusteeInput.Phone != nil {
		t.Phone = *trusteeInput.Phone
	}

	if trusteeInput.Email != nil {
		t.Email = *trusteeInput.Email
	}

	if trusteeInput.Relationship != nil {
		t.Relationship = trusteeInput.Relationship.String()
	}

	if trusteeInput.AdditionalInformation != nil {
		t.AdditionalInfo = nulls.NewString(*trusteeInput.AdditionalInformation)
	}

	if trusteeInput.TwitterLink != nil {
		t.TwitterLink = nulls.NewString(*trusteeInput.TwitterLink)
	}

	if trusteeInput.FacebookLink != nil {
		t.FacebookLink = nulls.NewString(*trusteeInput.FacebookLink)
	}

	verrs, err := t.Update(TX)

	if err != nil {
		return nil, errs.WithStack(err)
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

	return t.ToGraphQL(), nil
}

func (r *mutationResolver) DeleteTrustee(ctx context.Context, id scalars.UUID) (*model.Trustee, error) {
	trustee := &models.Trustee{}

	if err := TX.Find(trustee, id); err != nil {
		return nil, err
	}

	if err := TX.Destroy(trustee); err != nil {
		return nil, err
	}

	if trustee.UserID != User.ID {
		graphql.AddError(ctx, &gqlerror.Error{
			Message: "User does not own trustee",
		})
		return nil, nil
	}

	return trustee.ToGraphQL(), nil
}

func (r *queryResolver) Trustees(ctx context.Context, orderBy *string, order *string) ([]*model.Trustee, error) {
	var trustees models.Trustees

	nullableOrder := GetNullable(order)
	nullableOrderBy := GetNullable(orderBy)

	if err := trustees.GetForUser(TX, User.ID, nullableOrderBy.String, nullableOrder.String); err != nil {
		return nil, err
	}

	return trustees.ToGraphQL(), nil
}

func (r *queryResolver) Trustee(ctx context.Context, id scalars.UUID) (*model.Trustee, error) {
	trustee := &models.Trustee{}

	err := TX.Find(trustee, id)

	if err != nil {
		return nil, err
	}

	if trustee.UserID != User.ID {
		graphql.AddError(ctx, &gqlerror.Error{
			Message: "User does not own trustee",
		})
		return nil, nil
	}

	return trustee.ToGraphQL(), nil
}
