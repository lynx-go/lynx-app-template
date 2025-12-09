package session

import (
	"context"

	"github.com/lynx-go/lynx-app-template/pkg/bigid"
	"github.com/lynx-go/lynx-app-template/pkg/jsonapi"
)

func CurrentUser(ctx context.Context) bigid.ID {
	uid, _ := jsonapi.SessionUser(ctx)
	return bigid.ID(uid.Int64())
}
