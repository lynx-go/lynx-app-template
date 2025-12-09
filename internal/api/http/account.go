package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	apipb "github.com/lynx-go/lynx-app-template/genproto/api/v1"
	"github.com/lynx-go/lynx-app-template/internal/usecase"
	"github.com/lynx-go/lynx-app-template/pkg/errors"
)

type AccountAPI struct {
	uc *usecase.Account
}

func NewAccountAPI(
	uc *usecase.Account,
) *AccountAPI {
	return &AccountAPI{uc: uc}
}

func (api *AccountAPI) SignUp(ctx context.Context, req *apipb.SignUpRequest) (*apipb.SignUpResponse, error) {
	return api.uc.SignUp(ctx, req)
}

func (api *AccountAPI) Token(c echo.Context) error {
	ctx := c.Request().Context()
	grantType := c.QueryParam("grant_type")
	switch grantType {
	case "password":
		req := &apipb.TokenPasswordRequest{}
		if err := c.Bind(req); err != nil {
			return err
		}
		resp, err := api.uc.TokenByPassword(ctx, req)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, resp)
	default:
		return errors.New(http.StatusBadRequest, "invalid grant_type")
	}
}

func (api *AccountAPI) Logout(ctx echo.Context) error {

	return ctx.JSON(http.StatusOK, &apipb.TokenResponse{})
}
