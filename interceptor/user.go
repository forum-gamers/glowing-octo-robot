package interceptor

import (
	"context"

	"github.com/forum-gamers/glowing-octo-robot/pkg/user"
	"github.com/golang-jwt/jwt"
)

func (i *InterceptorImpl) GetUserFromCtx(ctx context.Context) user.User {
	var user user.User

	claim, ok := ctx.Value(CONTEXTUSERKEY).(jwt.MapClaims)
	if !ok {
		return user
	}

	for key, val := range claim {
		switch key {
		case "id":
			user.Id = val.(string)
		case "accountType":
			user.AccountType = val.(string)
		case "username":
			user.Username = val.(string)
		case "isVerified":
			user.IsVerified = val.(bool)
		default:
			continue
		}
	}
	return user
}
