package authorization_test

import (
	"context"
	"testing"
	"time"

	"myapp/internal/domain"
	"myapp/internal/service/authorization"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type mockUserGetter struct {
	getByNameFn func(ctx context.Context, username string) (domain.User, error)
}

func (m *mockUserGetter) GetByName(ctx context.Context, username string) (domain.User, error) {
	return m.getByNameFn(ctx, username)
}

type mockUserRolesGetter struct {
	listByUserIDFn func(ctx context.Context, userID int) ([]string, error)
}

func (m *mockUserRolesGetter) ListByUserID(ctx context.Context, userID int) ([]string, error) {
	return m.listByUserIDFn(ctx, userID)
}

type mockRefreshTokenRepo struct {
	setFn         func(ctx context.Context, refreshToken, userID string) error
	deleteFn      func(ctx context.Context, refreshToken string) error
	getAndDeleteFn func(ctx context.Context, refreshToken string) (string, error)
}

func (m *mockRefreshTokenRepo) Set(ctx context.Context, refreshToken, userID string) error {
	return m.setFn(ctx, refreshToken, userID)
}
func (m *mockRefreshTokenRepo) Delete(ctx context.Context, refreshToken string) error {
	return m.deleteFn(ctx, refreshToken)
}
func (m *mockRefreshTokenRepo) GetAndDelete(ctx context.Context, refreshToken string) (string, error) {
	return m.getAndDeleteFn(ctx, refreshToken)
}

type mockTokenBlacklist struct {
	setFn func(ctx context.Context, accessTokenHash string) error
}

func (m *mockTokenBlacklist) Set(ctx context.Context, accessTokenHash string) error {
	return m.setFn(ctx, accessTokenHash)
}

type mockJWTService struct {
	generateAccessFn  func(userID int, roles []string, expiresAt time.Duration) (string, error)
	generateRefreshFn func() (string, error)
	parseAccessFn     func(tokenString string) (authorization.AccessTokenClaims, error)
}

func (m *mockJWTService) GenerateAccessToken(userID int, roles []string, expiresAt time.Duration) (string, error) {
	return m.generateAccessFn(userID, roles, expiresAt)
}
func (m *mockJWTService) GenerateRefreshToken() (string, error) {
	return m.generateRefreshFn()
}
func (m *mockJWTService) ParseAccessToken(tokenString string) (authorization.AccessTokenClaims, error) {
	return m.parseAccessFn(tokenString)
}

func hashPassword(t *testing.T, password string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	require.NoError(t, err)
	return string(h)
}

func newAuthService(
	userGetter authorization.UserGetter,
	rolesGetter authorization.UserRolesGetter,
	refreshRepo authorization.RefreshTokenRepository,
	blacklist authorization.AccessTokenBlacklistSetter,
	jwt authorization.HMACJWTService,
) authorization.AuthService {
	return authorization.NewAuthService(userGetter, rolesGetter, refreshRepo, blacklist, jwt)
}

func TestAuthService_Auth(t *testing.T) {
	ctx := context.Background()
	password := "1234Qwerty"
	hash := hashPassword(t, password)

	t.Run("success", func(t *testing.T) {
		userGetter := &mockUserGetter{
			getByNameFn: func(ctx context.Context, username string) (domain.User, error) {
				assert.Equal(t, "alice", username)
				return domain.User{ID: 7, Username: "alice", PasswordHash: hash}, nil
			},
		}
		rolesGetter := &mockUserRolesGetter{
			listByUserIDFn: func(ctx context.Context, userID int) ([]string, error) {
				assert.Equal(t, 7, userID)
				return []string{"user"}, nil
			},
		}
		refreshRepo := &mockRefreshTokenRepo{
			setFn: func(ctx context.Context, refreshToken, userID string) error {
				assert.Equal(t, "7", userID)
				assert.NotEmpty(t, refreshToken)
				return nil
			},
		}
		jwtSvc := &mockJWTService{
			generateAccessFn: func(userID int, roles []string, expiresAt time.Duration) (string, error) {
				assert.Equal(t, 7, userID)
				assert.Equal(t, []string{"user"}, roles)
				return "access-token", nil
			},
			generateRefreshFn: func() (string, error) {
				return "refresh-token", nil
			},
		}

		svc := newAuthService(userGetter, rolesGetter, refreshRepo, &mockTokenBlacklist{}, jwtSvc)

		result, err := svc.Auth(ctx, authorization.LoginInput{Username: "alice", Password: password})
		require.NoError(t, err)
		assert.Equal(t, "access-token", result.AccessToken)
		assert.Equal(t, "refresh-token", result.RefreshToken)
		assert.Equal(t, []string{"user"}, result.Roles)
	})

	t.Run("user not found", func(t *testing.T) {
		userGetter := &mockUserGetter{
			getByNameFn: func(ctx context.Context, username string) (domain.User, error) {
				return domain.User{}, domain.ErrNotFound
			},
		}
		svc := newAuthService(userGetter, &mockUserRolesGetter{}, &mockRefreshTokenRepo{}, &mockTokenBlacklist{}, &mockJWTService{})

		_, err := svc.Auth(ctx, authorization.LoginInput{Username: "ghost", Password: password})
		require.ErrorIs(t, err, domain.ErrUnauthorized)
	})

	t.Run("wrong password", func(t *testing.T) {
		userGetter := &mockUserGetter{
			getByNameFn: func(ctx context.Context, username string) (domain.User, error) {
				return domain.User{ID: 7, Username: "alice", PasswordHash: hash}, nil
			},
		}
		svc := newAuthService(userGetter, &mockUserRolesGetter{}, &mockRefreshTokenRepo{}, &mockTokenBlacklist{}, &mockJWTService{})

		_, err := svc.Auth(ctx, authorization.LoginInput{Username: "alice", Password: "wrong"})
		require.ErrorIs(t, err, domain.ErrUnauthorized)
	})
}

func TestAuthService_Refresh(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		refreshRepo := &mockRefreshTokenRepo{
			getAndDeleteFn: func(ctx context.Context, refreshToken string) (string, error) {
				return "15", nil
			},
			setFn: func(ctx context.Context, refreshToken, userID string) error {
				assert.Equal(t, "15", userID)
				return nil
			},
		}
		rolesGetter := &mockUserRolesGetter{
			listByUserIDFn: func(ctx context.Context, userID int) ([]string, error) {
				assert.Equal(t, 15, userID)
				return []string{"admin"}, nil
			},
		}
		jwtSvc := &mockJWTService{
			generateAccessFn: func(userID int, roles []string, expiresAt time.Duration) (string, error) {
				return "new-access", nil
			},
			generateRefreshFn: func() (string, error) {
				return "new-refresh", nil
			},
		}

		svc := newAuthService(&mockUserGetter{}, rolesGetter, refreshRepo, &mockTokenBlacklist{}, jwtSvc)

		result, err := svc.Refresh(ctx, "old-refresh")
		require.NoError(t, err)
		assert.Equal(t, "new-access", result.AccessToken)
		assert.Equal(t, "new-refresh", result.RefreshToken)
		assert.Equal(t, []string{"admin"}, result.Roles)
	})

	t.Run("token not found", func(t *testing.T) {
		refreshRepo := &mockRefreshTokenRepo{
			getAndDeleteFn: func(ctx context.Context, refreshToken string) (string, error) {
				return "", domain.ErrNotFound
			},
		}
		svc := newAuthService(&mockUserGetter{}, &mockUserRolesGetter{}, refreshRepo, &mockTokenBlacklist{}, &mockJWTService{})

		_, err := svc.Refresh(ctx, "missing")
		require.ErrorIs(t, err, domain.ErrUnauthorized)
	})
}

func TestAuthService_Logout(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		refreshRepo := &mockRefreshTokenRepo{
			deleteFn: func(ctx context.Context, refreshToken string) error {
				assert.NotEmpty(t, refreshToken)
				return nil
			},
		}
		blacklist := &mockTokenBlacklist{
			setFn: func(ctx context.Context, accessTokenHash string) error {
				assert.NotEmpty(t, accessTokenHash)
				return nil
			},
		}

		svc := newAuthService(&mockUserGetter{}, &mockUserRolesGetter{}, refreshRepo, blacklist, &mockJWTService{})

		err := svc.Logout(ctx, "refresh", "access")
		require.NoError(t, err)
	})
}
