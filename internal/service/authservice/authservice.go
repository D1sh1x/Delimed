package authservice

import (
	"delimed/internal/repository"
	"delimed/internal/repository/models"
	"delimed/internal/transport/dto/request"
	"delimed/internal/utils/jwt"
	sl "delimed/internal/utils/logger"
	"delimed/internal/utils/password"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	RegisterUser(req request.SignUpInput) error
	LoginUser(req request.SignInInput) (string, error)
}

type AuthService struct {
	jwtSecret []byte
	repo      repository.RepositoryInterface
	logger    *slog.Logger
}

func NewAuthService(repo repository.RepositoryInterface, jwtSecret []byte, logger *slog.Logger) AuthServiceInterface {
	return &AuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *AuthService) RegisterUser(req request.SignUpInput) error {
	const op = "AuthService.RegisterUser"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("username", req.Username),
	)

	log.Info("attempting to register user")

	_, err := s.repo.GetUserByUsername(req.Username)
	if err == nil {
		log.Info("user already exists")
		return fmt.Errorf("user with username %s already exists", req.Username)
	}

	hashedPassword, err := password.HashPassword(req.Password)
	if err != nil {
		log.Error("error hashing password", sl.Err(err))
		return fmt.Errorf("error hashing password: %w", err)
	}

	user := &models.User{
		ID:       uuid.New(),
		Username: req.Username,
		Password: hashedPassword,
		Role:     "user",
	}
	if err := s.repo.CreateUser(user); err != nil {
		log.Error("error ccreating user in repository", sl.Err(err))
		return fmt.Errorf("error creating user in repository: %w", err)
	}

	return nil
}

func (s *AuthService) LoginUser(req request.SignInInput) (string, error) {
	const op = "AuthService.LoginUser"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("username", req.Username),
	)

	log.Info("attempting to login user")

	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		log.Info("failed to find user", sl.Err(err))
		return "", err
	}

	if err := password.VerifyPassword(user.Password, req.Password); err != nil {
		log.Info("invalid email or password", sl.Err(err))
		return "", fmt.Errorf("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, s.jwtSecret)
	if err != nil {
		log.Warn("error generating token", sl.Err(err))
		return "", fmt.Errorf("error generating token: %w", err)
	}

	return token, nil
}
