package directives

import (
	"context"
	"log"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/helpers"
	"github.com/khoirulhasin/untirta_api/app/infrastructures/middlewares"
	"github.com/khoirulhasin/untirta_api/app/models"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func HasRoleDirective(ctx context.Context, obj interface{}, next graphql.Resolver, roles []models.RoleEnum) (interface{}, error) {
	user := middlewares.ForContext(ctx)
	if user == nil {
		return nil, gqlerror.Errorf("Tidak diizinkan: Harap login")
	}

	// Convert input roles to []string and validate
	requiredRoles := make([]string, 0, len(roles))
	for _, rolePtr := range roles {
		if rolePtr != "" {
			normalizedRole := strings.ToLower(string(rolePtr)) // Convert RoleEnum to string and lowercase
			requiredRoles = append(requiredRoles, normalizedRole)

			if !helpers.Contains([]string{"admin", "operator", "user", "driver"}, normalizedRole) {
				log.Printf("Invalid role: %s", rolePtr)
				return nil, gqlerror.Errorf("Peran tidak valid: %s", rolePtr)
			}
			requiredRoles = append(requiredRoles, normalizedRole)
		}
	}

	// Check if user has at least one required role
	for _, requiredRole := range requiredRoles {
		if helpers.Contains(user.Roles, requiredRole) {
			return next(ctx)
		}
	}

	// Log failed role check for debugging
	log.Printf("User %d (%s) lacks required roles: %v (user roles: %v)", user.ID, user.Name, requiredRoles, user.Roles)
	return nil, gqlerror.Errorf("Izin tidak cukup: Diperlukan salah satu peran %v", requiredRoles)
}
