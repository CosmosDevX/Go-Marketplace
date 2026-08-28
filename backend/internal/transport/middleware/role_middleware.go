package middleware

import (
	"fmt"
	"myapp/internal/domain"
	"myapp/internal/utils"
	"net/http"
	"slices"
)

func RoleMiddleware(rolesToAccess []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user, err := UserFromContext(ctx)
			if err != nil {
				utils.WriteError(ctx, w, err)
				return
			}

			for i := range user.Roles {
				if slices.Contains(rolesToAccess, user.Roles[i]) {
					next.ServeHTTP(w, r)
					return
				}
			}

			utils.WriteError(ctx, w, fmt.Errorf("access denied: %w", domain.ErrForbidden))
		})
	}
}
