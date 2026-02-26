package directives

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/directives/validates"
)

func ValidateDirective(ctx context.Context, obj interface{}, next graphql.Resolver,
	required *bool,
	email *bool,
	username *bool,
	password *bool,
	integer *bool) (interface{}, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}

	// Get field context for error reporting
	fieldCtx := graphql.GetFieldContext(ctx)
	fieldName := fieldCtx.Field.Name
	if fieldName == "" {
		fieldName = "unknown"
	}

	// Skip validation if required is nil or false
	if required != nil && *required {
		if _, err := validates.ValidateRequired(val, fieldName); err != nil {
			return nil, err
		}
	}

	// Skip validation if email is nil or false
	if email != nil && *email {
		if _, err := validates.ValidateEmail(val, fieldName); err != nil {
			return nil, err
		}
	}

	// Skip validation if username is nil or false
	if username != nil && *username {
		if _, err := validates.ValidateUsername(val, fieldName); err != nil {
			return nil, err
		}
	}

	if password != nil && *password {
		if _, err := validates.ValidatePassword(val, fieldName); err != nil {
			return nil, err
		}
	}

	if integer != nil && *integer {
		if _, err := validates.ValidateInteger(val, fieldName); err != nil {
			return nil, err
		}
	}

	return val, nil
}
