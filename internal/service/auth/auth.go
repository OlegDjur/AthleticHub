package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"time"
	"workout/internal/adapters/postgres"
	"workout/internal/dto"
	"workout/internal/lib/jwt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	auth     Auth
	tokenTTL time.Duration
	secret   string
}

func NewAuthService(auth Auth, ttl time.Duration) *AuthService {
	return &AuthService{
		auth:     auth,
		tokenTTL: ttl,
	}
}

func (a *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	const op = "Auth.Login"

	//log := a.log.With(
	//	slog.String("op", op),
	//	slog.String("email", email))
	//
	//log.Info("attempting to login user")

	user, err := a.auth.GetUser(ctx, req.Login)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			//a.log.Warn("user not found")

			return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		//a.log.Error("failed to get user", err)

		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		//a.log.Info("invalid credentials")

		return nil, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	//log.Info("user logged in successfully")
	fmt.Println(a.tokenTTL)
	token, err := jwt.NewToken(user, a.secret, a.tokenTTL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &dto.LoginResponse{Token: token}, nil
}

func (a *AuthService) RegisterNewUser(ctx context.Context, req dto.RegisterRequest) (string, error) {
	const op = "auth.RegisterNewUser"

	//log := a.log.With(
	//	slog.String("op", op),
	//	slog.String("email", email))
	//
	//log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash:", err.Error())
		return "", err
	}

	id, err := a.auth.CreateUser(ctx, req.Login, passHash)
	if err != nil {
		if errors.Is(err, postgres.ErrUserExists) {
			log.Warn("user already exists")

			return "", fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user:", err.Error())

		return "", err
	}

	log.Info("user created:", id)

	return id, nil

}

func (a *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	//log := a.log.With(
	//	slog.String("op", op),
	//	slog.String("user_id", fmt.Sprint(userID)))
	//
	//log.Info("checking if user is admin")

	isAdmin, err := a.auth.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
	}

	log.Info("checking if user is admin")

	return isAdmin, nil
}
