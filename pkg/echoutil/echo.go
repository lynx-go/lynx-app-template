package echoutil

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lynx-go/lynx-app-template/pkg/errors"
	"github.com/lynx-go/lynx-app-template/pkg/jsonapi"
	"github.com/lynx-go/x/log"
	"github.com/spf13/cast"
)

type Options struct {
	JWTSecret string `json:"jwt_secret"`
}

func NewEcho(logger *slog.Logger) *echo.Echo {

	router := echo.New()
	router.Use(middleware.Recover())
	router.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		RequestIDHandler: func(e echo.Context, requestId string) {
			ctx := e.Request().Context()
			reqLogger := log.FromContext(ctx).With("x-request-id", requestId)
			ctx = log.Context(ctx, reqLogger)
			e.SetRequest(e.Request().WithContext(ctx))
		},
		TargetHeader: "X-Request-ID",
	}))
	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			args := []any{
				"request_id", v.RequestID,
				"method", v.Method,
				"url", v.URI,
				"status", v.Status,
				"latency", v.Latency,
				"remote_ip", v.RemoteIP,
				"host", v.Host,
				"user_agent", v.UserAgent,
			}
			if v.Error != nil {
				args = append(args, "error", v.Error.Error())
			}
			logger.DebugContext(c.Request().Context(), "[echo] access log", args...)
			return nil
		},
		LogLatency:   true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogURIPath:   true,
		LogRoutePath: true,
		LogRequestID: true,
		LogStatus:    true,
		LogError:     true,
	}))

	router.HTTPErrorHandler = func(err error, c echo.Context) {
		var apiErr *errors.APIError
		switch e := err.(type) {
		case *errors.APIError:
			if e.Code == 401 {
				_ = c.JSON(http.StatusUnauthorized, &errResponse{Error: e})
				return
			}
			apiErr = e
		case *echo.HTTPError:
			apiErr = &errors.APIError{
				Code:    e.Code,
				Message: cast.ToString(e.Message),
				Details: map[string]interface{}{
					"internal": e.Error(),
				},
			}
		default:
			apiErr = &errors.APIError{
				Code:    -1,
				Message: "服务异常，请稍后再试",
				Details: map[string]any{
					"internal": err.Error(),
				},
			}
		}
		_ = c.JSON(200, &errResponse{Error: apiErr})
	}
	return router
}

func JWTMiddleware(jwtSecret string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		ContinueOnIgnoredError: false,
		SigningKey:             []byte(jwtSecret),
		ContextKey:             jsonapi.ContextKeyJWT,
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &SessionTokenClaims{}
		},
		ErrorHandler: func(c echo.Context, err error) error {
			newErr := errors.New(401, "认证失败")
			return newErr.Wrap(err)
		},
	})
}

type errResponse struct {
	Error *errors.APIError `json:"error"`
}

type SessionTokenClaims struct {
	TokenId   string            `json:"tid,omitempty"`
	UserId    string            `json:"uid,omitempty"`
	Username  string            `json:"usn,omitempty"`
	Vars      map[string]string `json:"vrs,omitempty"`
	ExpiresAt int64             `json:"exp,omitempty"`
	IssuedAt  int64             `json:"iat,omitempty"`
}

func (s *SessionTokenClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(s.ExpiresAt, 0)), nil
}
func (s *SessionTokenClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}
func (s *SessionTokenClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(s.IssuedAt, 0)), nil
}
func (s *SessionTokenClaims) GetAudience() (jwt.ClaimStrings, error) {
	return []string{}, nil
}
func (s *SessionTokenClaims) GetIssuer() (string, error) {
	return "", nil
}
func (s *SessionTokenClaims) GetSubject() (string, error) {
	return s.UserId, nil
}
