package account

import (
	"context"
	"time"

	"github.com/lynx-go/lynx-app-template/pkg/bigid"
)

type UsersRepo interface {
	Get(ctx context.Context, id bigid.ID) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, val UserCreate) error
}

type UserCreate struct {
	ID           bigid.ID  `json:"id"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	CreatedBy    bigid.ID  `json:"created_by"`
}

type User struct {
	ID               bigid.ID               `json:"id"`
	Username         string                 `json:"username"`
	DisplayName      string                 `json:"display_name"`
	PasswordHash     string                 `json:"password_hash"`
	AvatarURL        string                 `json:"avatar_url"`
	Phone            string                 `json:"phone"`
	PhoneConfirmedAt time.Time              `json:"phone_confirmed_at"`
	Email            string                 `json:"email"`
	EmailConfirmedAt time.Time              `json:"email_confirmed_at"`
	Status           int8                   `json:"status"`
	Gender           int8                   `json:"gender"`
	ConfirmedAt      time.Time              `json:"confirmed_at"`
	Role             string                 `json:"role"`
	AppMetadata      map[string]interface{} `json:"app_metadata"`
	UserMetadata     map[string]interface{} `json:"user_metadata"`
	LastSignInAt     time.Time              `json:"last_sign_in_at"`
	BannedUntil      time.Time              `json:"banned_until"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}
