package repoimpl

import (
	"context"

	"github.com/lynx-go/lynx-app-template/internal/domain/account"
	"github.com/lynx-go/lynx-app-template/internal/infra/clients"
	entgen "github.com/lynx-go/lynx-app-template/internal/infra/ent/gen"
	"github.com/lynx-go/lynx-app-template/internal/infra/ent/gen/user"
	"github.com/lynx-go/lynx-app-template/pkg/bigid"
)

type usersImpl struct {
	*clients.DataClients
}

func (impl *usersImpl) Get(ctx context.Context, id bigid.ID) (*account.User, error) {
	v, err := impl.DB.User(ctx).Get(ctx, id.Int64())
	if entgen.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return newUser(v), nil
}

func (impl *usersImpl) Create(ctx context.Context, v account.UserCreate) error {
	return impl.DB.User(ctx).Create().
		SetID(v.ID.Int64()).
		SetUsername(v.Username).
		SetDisplayName(v.DisplayName).
		SetEmail(v.Email).
		SetRole(v.Role).
		SetPasswordHash(v.Password).
		SetCreatedAt(v.CreatedAt).
		SetUpdatedAt(v.CreatedAt).
		SetIsSuperAdmin(v.IsSuperAdmin).
		SetCreatedBy(v.CreatedBy.Int64()).
		SetUpdatedBy(v.CreatedBy.Int64()).
		Exec(ctx)
}

func (impl *usersImpl) GetByEmail(ctx context.Context, email string) (*account.User, error) {
	v, err := impl.DB.User(ctx).Query().Where(user.Email(email), user.StatusNEQ(-2)).First(ctx)
	if entgen.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return newUser(v), nil
}

func newUser(v *entgen.User) *account.User {
	return &account.User{
		ID:               bigid.ID(v.ID),
		Username:         v.Username,
		DisplayName:      v.DisplayName,
		PasswordHash:     v.PasswordHash,
		AvatarURL:        v.AvatarURL,
		Phone:            v.Phone,
		PhoneConfirmedAt: v.PhoneConfirmedAt,
		Email:            v.Email,
		EmailConfirmedAt: v.EmailConfirmedAt,
		Status:           v.Status,
		Gender:           v.Gender,
		ConfirmedAt:      v.ConfirmedAt,
		Role:             v.Role,
		AppMetadata:      v.AppMetadata,
		UserMetadata:     v.UserMetadata,
		LastSignInAt:     v.LastSignInAt,
		BannedUntil:      v.BannedUntil,
		CreatedAt:        v.CreatedAt,
		UpdatedAt:        v.UpdatedAt,
	}
}

func (impl *usersImpl) GetByUsername(ctx context.Context, username string) (*account.User, error) {
	v, err := impl.DB.User(ctx).Query().Where(user.Username(username), user.StatusNEQ(-2)).First(ctx)
	if entgen.IsNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return newUser(v), nil
}

func NewUserRepo(
	data *clients.DataClients,
) account.UsersRepo {
	return &usersImpl{DataClients: data}
}
