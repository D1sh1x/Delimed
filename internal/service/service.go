package service

import (
	"delimed/internal/repository"
	"delimed/internal/service/authservice"
	"delimed/internal/service/deliveryservice"
	"delimed/internal/service/userservice"
	"log/slog"
)

type Service struct {
	auth     authservice.AuthServiceInterface
	user     userservice.UserServiceInterface
	delivery deliveryservice.DeliveryServiceInterface
}

type ServiceInterface interface {
	Auth() authservice.AuthServiceInterface
	User() userservice.UserServiceInterface
	Delivery() deliveryservice.DeliveryServiceInterface
}

func (s *Service) Auth() authservice.AuthServiceInterface             { return s.auth }
func (s *Service) User() userservice.UserServiceInterface             { return s.user }
func (s *Service) Delivery() deliveryservice.DeliveryServiceInterface { return s.delivery }

func NewService(repo repository.RepositoryInterface, jwtSecret []byte, logger *slog.Logger, cdekClientID, cdekClientSecret, dellinAppKey string) *Service {

	return &Service{
		auth:     authservice.NewAuthService(repo, jwtSecret, logger),
		user:     userservice.NewUserService(repo, logger),
		delivery: deliveryservice.NewDeliveryService(logger, cdekClientID, cdekClientSecret, dellinAppKey),
	}
}
