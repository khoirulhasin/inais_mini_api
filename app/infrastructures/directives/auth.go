package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/middlewares"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func AuthDirective(ctx context.Context, obj any, next graphql.Resolver) (res any, err error) {
	user := middlewares.ForContext(ctx)
	if user == nil {
		return nil, gqlerror.Errorf("Tidak diizinkan: Harap login")
	}
	return next(ctx)
}
