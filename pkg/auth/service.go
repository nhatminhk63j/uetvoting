package auth

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/nhatminhk63j/uetvoting/config"
	"github.com/nhatminhk63j/uetvoting/pkg/jwt"
	"github.com/nhatminhk63j/uetvoting/pkg/logger"
	"github.com/nhatminhk63j/uetvoting/pkg/user"
)

const (
	googleOauthUrl = "https://oauth2.googleapis.com/tokeninfo?id_token="
)

// Service ...
type Service interface {
	Login(ctx context.Context, idToken string) (serverToken string, err error)
}

// UserService ...
type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
	UpsertUser(ctx context.Context, user *user.User) error
}

type service struct {
	jwtManager   *jwt.Resolver
	httpClient   *http.Client
	oathClientID string

	userSvc UserService
}

// NewService ...
func NewService(jwtManager *jwt.Resolver, userSvc UserService) Service {
	authCfg := config.LoadAuthConfig()

	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &service{
		jwtManager:   jwtManager,
		userSvc:      userSvc,
		oathClientID: authCfg.OauthClientID,
		httpClient:   &http.Client{Transport: httpTransport},
	}
}

// Login ...
func (s service) Login(ctx context.Context, idToken string) (serverToken string, err error) {
	userGoogle, err := s.GetUserByGoogleAccessToken(idToken)
	if err != nil {
		err = xerrors.Errorf("error get user by google id token: %w", err)
		return
	}

	userInfo, err := s.userSvc.GetUserByEmail(ctx, userGoogle.Email)
	if err != nil {
		var notFoundErr *user.ErrNotFound
		if !errors.As(err, &notFoundErr) {
			err = xerrors.Errorf("error getting user when login: %w", err)
			return "", err
		}
		newUser := &user.User{
			Email:   userGoogle.Email,
			Name:    userGoogle.Name,
			Picture: userGoogle.Picture,
			Role:    UserRole,
		}
		err := s.userSvc.UpsertUser(ctx, newUser)
		if err != nil {
			return "", xerrors.Errorf("error create new user: %w", err)
		}
		return s.jwtManager.GenerateToken(newUser)
	}

	userInfo.Name = userGoogle.Name
	userInfo.Picture = userGoogle.Picture
	err = s.userSvc.UpsertUser(ctx, userInfo)
	if err != nil {
		return "", xerrors.Errorf("error update user: %w", err)
	}
	return s.jwtManager.GenerateToken(userInfo)
}

// GetUserByGoogleAccessToken get user by Google AccessToken.
func (s service) GetUserByGoogleAccessToken(idToken string) (*TokenInfoGoogle, error) {
	// Get token info from Google oauth2 with idToken.
	resp, err := s.httpClient.Get(googleOauthUrl + idToken)
	if err != nil {
		logger.Errorf("error get response from google oauth server: %v", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body) // Read resp.Body

	if err != nil {
		return nil, err
	}

	tokenInfo, err := s.handleErrorResponseFromGoogle(body)
	if err != nil {
		return nil, err
	}
	if !s.hasPermissionAccessApp(tokenInfo) {
		err = NoPermissionError{}
		return nil, err
	}
	return tokenInfo, nil
}

// handleErrorResponseFromGoogle handle error response from Google Auth.
func (s service) handleErrorResponseFromGoogle(body []byte) (*TokenInfoGoogle, error) {
	var errorGoogle ErrorGoogle
	err := json.Unmarshal(body, &errorGoogle)
	if err != nil {
		return nil, xerrors.Errorf("err unmarshal body to error google: %w", err)
	}

	if errorGoogle.Error.Code == 401 {
		return nil, xerrors.Errorf("unauthenticated by google id token")
	}

	var userGoogle TokenInfoGoogle
	err = json.Unmarshal(body, &userGoogle) // Parse []byte code response from google to token info object.
	if err != nil {
		return nil, xerrors.Errorf("error unmarshal body to token info object: %w", err)
	}
	return &userGoogle, nil
}

// hasPermissionAccessApp check app id of token info.
func (s service) hasPermissionAccessApp(tokenInfo *TokenInfoGoogle) bool {
	return tokenInfo.Aud == s.oathClientID
}
