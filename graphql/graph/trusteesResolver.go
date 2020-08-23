package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/kolioDev/after_life/graphql/graph/model"
	"github.com/kolioDev/after_life/models"
	"github.com/pkg/errors"
	"log"
	"time"
)

func (r *mutationResolver) CreateTrustee(ctx context.Context, input model.NewTrustee) (*model.Trustee, error) {

	t := models.Trustee{
		Name:         input.Name,
		Relationship: input.Relationship.String(),
		Phone:        input.Phone,
		Email:        input.Email,
		FacebookLink: GetNullable(input.FacebookLink),
		//TwitterLink:    nulls.NewString(*input.TwitterLink),
		//AdditionalInfo: nulls.NewString(*input.AdditionalInformation),
	}

	fmt.Println("INPUT EMAIL", input.Email)
	fmt.Println("TRUSTEE EMAIL", t.Email)

	verrs, err := t.Create(models.DB, User.ID)
	if err != nil {
		return nil, err
	}
	if verrs.HasAny() {
		jsonErrs, _ := json.Marshal(verrs.Errors)
		return nil, gqlerrors.FormatError(errors.New(string(jsonErrs)))
	}

	return t.ToGraphQL(), nil
}

func (r *queryResolver) Trustees(ctx context.Context) ([]*model.Trustee, error) {
	var trustees []*model.Trustee

	log.Println(User)

	trustees = append(trustees, models.Trustee{
		ID:             uuid.UUID{},
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Name:           "Go6o",
		Relationship:   "friend",
		Phone:          "0882726289",
		Email:          "nikitvmaniak@gmail.com",
		FacebookLink:   nulls.String{},
		TwitterLink:    nulls.String{},
		AdditionalInfo: nulls.String{},
	}.ToGraphQL())

	return trustees, nil
}
