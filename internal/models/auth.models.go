package models

import "context"

type UserContext struct {
	Id   string
	Name string
}

func SetContextUser(ctx context.Context, u *UserContext) context.Context {
	return context.WithValue(ctx, "user", u)
}

func GetContextUser(ctx context.Context) *UserContext {

	user, ok := ctx.Value("user").(*UserContext)

	if !ok {
		return nil
	}

	return user

}
