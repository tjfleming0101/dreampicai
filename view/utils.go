package view

import (
	"context"

	"github.com/tjfleming0101/dreampicai/types"
)

// can pass authenticatedUser to Navigation
func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value(types.UserContextKey).(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}