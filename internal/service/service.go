package service

import (
	"delimed/internal/repository"
	"delimed/internal/service/authservice"
	"delimed/internal/service/userservice"
	"log/slog"
)

type Service struct {
	auth authservice.AuthServiceInterface
	user userservice.UserServiceInterface
}

type ServiceInterface interface {
	Auth() authservice.AuthServiceInterface
	User() userservice.UserServiceInterface
}

func (s *Service) Auth() authservice.AuthServiceInterface { return s.auth }
func (s *Service) User() userservice.UserServiceInterface { return s.user }

func NewService(repo repository.RepositoryInterface, jwtSecret []byte, logger *slog.Logger) *Service {

	return &Service{
		auth: authservice.NewAuthService(repo, jwtSecret, logger),
		user: userservice.NewUserService(repo, logger),
	}
}
