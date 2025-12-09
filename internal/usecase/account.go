package usecase

import (
	"context"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	apipb "github.com/lynx-go/lynx-app-template/genproto/api/v1"
	"github.com/lynx-go/lynx-app-template/internal/domain/account"
	"github.com/lynx-go/lynx-app-template/internal/domain/runtimevars"
	configpb "github.com/lynx-go/lynx-app-template/internal/pkg/config"
	"github.com/lynx-go/lynx-app-template/internal/pkg/varkeys"
	"github.com/lynx-go/lynx-app-template/pkg/bigid"
	"github.com/lynx-go/lynx-app-template/pkg/echoutil"
	"github.com/lynx-go/x/log"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	repo        account.UsersRepo
	ids         bigid.IDGen
	config      *configpb.AppConfig
	runtimeVars runtimevars.RuntimeVars
}

var (
	invalidUsernameRegex = regexp.MustCompilePOSIX("([[:cntrl:]]|[[\t\n\r\f\v]])+")
	invalidCharsRegex    = regexp.MustCompilePOSIX("([[:cntrl:]]|[[:space:]])+")
	emailRegex           = regexp.MustCompile(`^.+@.+\..+$`)
)

func NewAccount(
	repo account.UsersRepo,
	config *configpb.AppConfig,
	runtimeVars runtimevars.RuntimeVars,
) *Account {
	return &Account{
		repo:        repo,
		ids:         bigid.NewIDGen(),
		config:      config,
		runtimeVars: runtimeVars,
	}
}

func (uc *Account) TokenByPassword(ctx context.Context, req *apipb.TokenPasswordRequest) (*apipb.TokenResponse, error) {
	email := req.GetEmail()
	username := req.GetUsername()
	password := req.GetPassword()
	if len(password) == 0 {
		return nil, errors.New("用户名或密码错误")
	}
	var user *account.User
	var err error
	if username != "" {
		user, err = uc.repo.GetByUsername(ctx, username)
		if err != nil {
			return nil, errors.Wrap(err, "查询用户信息失败")
		}
		if user == nil {
			return nil, errors.New("用户名或密码错误")
		}
	} else if email != "" {
		user, err = uc.repo.GetByEmail(ctx, email)
		if err != nil {
			return nil, errors.Wrap(err, "查询用户信息失败")
		}
		if user == nil {
			return nil, errors.New("Email或密码错误")
		}
	}
	if user == nil {
		return nil, errors.New("Email或密码错误")
	}
	if len(user.PasswordHash) == 0 {
		return nil, errors.New("用户名或密码错误")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	tokenID := uuid.NewString()
	tokenIssuedAt := time.Now().Unix()
	meta := map[string]string{}
	token, exp := generateToken(uc.config, tokenID, tokenIssuedAt, user.ID.String(), username, meta)
	refreshToken, refreshExp := generateRefreshToken(uc.config, tokenID, tokenIssuedAt, user.ID.String(), username, meta)
	return &apipb.TokenResponse{
		Token:                 token,
		ExpiresAt:             exp,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExp,
		UserInfo:              newSessionUser(user),
	}, nil
}

func generateToken(config *configpb.AppConfig, tokenID string, tokenIssuedAt int64, userID, username string, vars map[string]string) (string, int64) {
	exp := time.Now().UTC().Add(time.Duration(config.GetSecurity().GetJwt().TokenExpirySec) * time.Second).Unix()
	return generateTokenWithExpiry(config.GetSecurity().GetJwt().GetSecret(), tokenID, tokenIssuedAt, userID, username, vars, exp)
}

func generateRefreshToken(config *configpb.AppConfig, tokenID string, tokenIssuedAt int64, userID string, username string, vars map[string]string) (string, int64) {
	exp := time.Now().UTC().Add(time.Duration(config.GetSecurity().GetJwt().RefreshTokenExpirySec) * time.Second).Unix()
	return generateTokenWithExpiry(config.GetSecurity().GetJwt().GetRefreshTokenSecret(), tokenID, tokenIssuedAt, userID, username, vars, exp)
}

func generateTokenWithExpiry(signingKey, tokenID string, tokenIssuedAt int64, userID, username string, vars map[string]string, exp int64) (string, int64) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &echoutil.SessionTokenClaims{
		TokenId:   tokenID,
		UserId:    userID,
		Username:  username,
		Vars:      vars,
		ExpiresAt: exp,
		IssuedAt:  tokenIssuedAt,
	})
	signedToken, _ := token.SignedString([]byte(signingKey))
	return signedToken, exp
}
func (uc *Account) SignUp(ctx context.Context, req *apipb.SignUpRequest) (*apipb.SignUpResponse, error) {
	email := req.GetEmail()
	username := req.GetUsername()
	user, err := uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "查询用户信息失败")
	}
	if user != nil {
		return nil, errors.New("用户名已注册")
	}
	uid, err := uc.ids.NextID()
	if err != nil {
		return nil, errors.Wrap(err, "生成用户ID失败")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "生成用户密码失败")
	}
	initialized, err := uc.runtimeVars.GetBool(ctx, varkeys.AppInitialized)
	if err != nil {
		log.ErrorContext(ctx, "query app initialized error", err)
	}
	isSuperAdmin := !initialized
	if err := uc.repo.Create(ctx, account.UserCreate{
		ID:           uid,
		Username:     username,
		DisplayName:  username,
		Email:        email,
		Role:         "normal",
		Password:     string(passwordHash),
		CreatedAt:    time.Now(),
		IsSuperAdmin: isSuperAdmin,
	}); err != nil {
		return nil, errors.Wrap(err, "生成用户失败")
	}

	user, err = uc.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "查询用户出错")
	}
	if err := uc.runtimeVars.Set(ctx, varkeys.AppInitialized, true); err != nil {
		log.ErrorContext(ctx, "update app initialized error", err)
	}
	return &apipb.SignUpResponse{
		UserInfo: newSessionUser(user),
	}, nil
}

func newSessionUser(v *account.User) *apipb.UserInfo {
	return &apipb.UserInfo{
		Id:          v.ID.Int64(),
		Username:    v.Username,
		DisplayName: v.DisplayName,
		AvatarUrl:   v.AvatarURL,
		Email:       v.Email,
		Phone:       v.Phone,
	}
}
