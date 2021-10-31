package jwt

import (
	"github.com/nhatminhk63j/uetvoting/pkg/user"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/xerrors"

	"github.com/nhatminhk63j/uetvoting/config"
)

type Resolver struct {
	secretKey     string
	TokenDuration time.Duration
}

type UserClaims struct {
	jwt.StandardClaims
	ID    int    `json:"id"`
	Email string `json:"email"`
	Role  int    `json:"role"`
}

// NewJWTResolver ...
func NewJWTResolver() *Resolver {
	cfg := config.LoadJWTConfig()
	return &Resolver{
		secretKey:     cfg.SecretKey,
		TokenDuration: time.Duration(cfg.TimeDuration) * time.Hour,
	}
}

// GenerateToken ...
func (j *Resolver) GenerateToken(user *user.User) (accessToken string, err error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.TokenDuration).Unix(),
		},
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err = token.SignedString([]byte(j.secretKey))
	if err != nil {
		err = xerrors.Errorf("failed to sign access token with secret key: %w", err)
		return
	}
	return
}

// Verify accessToken.
func (j *Resolver) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, xerrors.Errorf("unexpected token signing method")
			}
			return []byte(j.secretKey), nil
		},
	)
	if err != nil {
		return nil, xerrors.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)

	// check expiresAt of access token.
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, xerrors.Errorf("access token has expired")
	}
	if !ok {
		return nil, xerrors.Errorf("invalid token claims")
	}

	return claims, nil
}
