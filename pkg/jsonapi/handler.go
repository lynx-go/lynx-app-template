package jsonapi

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
)

type HandlerFunc[I any, O any] func(ctx context.Context, in I) (O, error)

func H[I any, O any](h HandlerFunc[*I, *O]) echo.HandlerFunc {
	return func(ec echo.Context) error {
		in := new(I)
		if err := ec.Bind(in); err != nil {
			return err
		}
		token, ok := ec.Get(ContextKeyJWT).(*jwt.Token)
		ctx := ec.Request().Context()
		if ok {
			sub, _ := token.Claims.GetSubject()
			ctx = context.WithValue(ctx, ContextKeyCurrentUser, sub)
			ec.SetRequest(ec.Request().WithContext(ctx))
		}
		out, err := h(ctx, in)
		if err != nil {
			return err
		}

		return ec.JSON(http.StatusOK, out)
	}
}

func SessionUser(ctx context.Context) (UserID, bool) {
	v, ok := ctx.Value(ContextKeyCurrentUser).(string)
	return UserID(v), ok
}

type UserID string

func (uid UserID) String() string {
	return string(uid)
}

func (uid UserID) Int64() int64 {
	return cast.ToInt64(uid.String())
}

func (uid UserID) Valid() bool {
	return uid != ""
}

const ContextKeyJWT = "lx_jwt"

const ContextKeyCurrentUser = "lx_user"
