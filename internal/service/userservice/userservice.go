package userservice

import (
	"bytes"
	"context"
	"delimed/internal/repository"
	"delimed/internal/transport/dto/request"
	"delimed/internal/transport/dto/response"
	"delimed/internal/utils/cdek"
	sl "delimed/internal/utils/logger"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type UserServiceInterface interface {
	GetTariffsList(ctx context.Context, cdekreq request.CDEKRequestList) (*response.CDEKTariffListResponse, error)
	GetTarifs(ctx context.Context, cdekreq request.CDEKRequest) (*response.CDEKTariffCalcResponse, error)
	GetUserByID(userID string) (*response.UserResponse, error)
	DeleteUser(userID string) error
}

type UserService struct {
	repo   repository.RepositoryInterface
	logger *slog.Logger
}

func NewUserService(repo repository.RepositoryInterface, logger *slog.Logger) UserServiceInterface {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

func (s *UserService) GetTariffsList(ctx context.Context, cdekreq request.CDEKRequestList) (*response.CDEKTariffListResponse, error) {
	const op = "UserService.GetTariffsList"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("from", cdekreq.From),
		slog.String("to", cdekreq.To),
		slog.Int("weight", cdekreq.Weight),
		slog.Int("length", cdekreq.Length),
		slog.Int("width", cdekreq.Width),
		slog.Int("height", cdekreq.Height),
	)

	log.Info("attempting to get CDEK tariffs list")

	const baseURL = "https://api.edu.cdek.ru/v2"
	const clientID = "wqGwiQx0gg8mLtiEKsUinjVSICCjtTEP"
	const clientSecret = "RmAmgvSgSl1yirlz9QupbzOJVqhCxcP5"

	token, err := cdek.GetCDEKToken(ctx, clientID, clientSecret)
	if err != nil {
		log.Error("failed to get CDEK token", sl.Err(err))
		return nil, fmt.Errorf("get token: %w", err)
	}

	log.Debug("CDEK token obtained successfully")

	pkg := request.CDEKPackage{
		Weight: cdekreq.Weight,
		Length: cdekreq.Length,
		Width:  cdekreq.Width,
		Height: cdekreq.Height,
	}

	reqBody := request.CDEKTariffCalcRequestList{
		Type:     1,
		Date:     time.Now().Format("2006-01-02T15:04:05-0700"),
		Currency: 1,
		Lang:     "rus",
		From:     request.CDEKLocation{Address: cdekreq.From},
		To:       request.CDEKLocation{Address: cdekreq.To},
		Packages: []request.CDEKPackage{pkg},
		// Services: []request.CDEKService{...} — если понадобятся
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("failed to marshal request body", sl.Err(err))
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/calculator/tarifflist",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		log.Error("failed to build HTTP request", sl.Err(err))
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	log.Debug("sending request to CDEK API", slog.String("url", baseURL+"/calculator/tarifflist"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("failed to execute HTTP request", sl.Err(err))
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	log.Debug("received response from CDEK API", slog.Int("status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		bodyStr := string(b)

		log.Error("tariff list calculation failed",
			slog.Int("status_code", resp.StatusCode),
			slog.String("status", resp.Status),
			slog.String("response_body", bodyStr),
		)
		return nil, fmt.Errorf("tarifflist calc failed: status=%s body=%s", resp.Status, bodyStr)
	}

	var tariffs response.CDEKTariffListResponse
	if err := json.NewDecoder(resp.Body).Decode(&tariffs); err != nil {
		log.Error("failed to decode response", sl.Err(err))
		return nil, fmt.Errorf("decode response: %w", err)
	}

	log.Info("successfully retrieved CDEK tariffs list")

	return &tariffs, nil
}

func (s *UserService) GetTarifs(ctx context.Context, cdekreq request.CDEKRequest) (*response.CDEKTariffCalcResponse, error) {
	const op = "UserService.GetTarifs"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("from", cdekreq.From),
		slog.String("to", cdekreq.To),
		slog.Int("tariff_code", cdekreq.TariffCode),
		slog.Int("weight", cdekreq.Weight),
	)

	log.Info("attempting to get CDEK tariffs")

	const baseURL = "https://api.edu.cdek.ru/v2"
	const clientID = "wqGwiQx0gg8mLtiEKsUinjVSICCjtTEP"
	const clientSecret = "RmAmgvSgSl1yirlz9QupbzOJVqhCxcP5"

	token, err := cdek.GetCDEKToken(ctx, clientID, clientSecret)
	if err != nil {
		log.Error("failed to get CDEK token", sl.Err(err))
		return nil, fmt.Errorf("get token: %w", err)
	}

	log.Debug("CDEK token obtained successfully")

	pkg := request.CDEKPackage{
		Weight: cdekreq.Weight,
		Length: cdekreq.Length,
		Width:  cdekreq.Width,
		Height: cdekreq.Height,
	}

	reqBody := request.CDEKTariffCalcRequest{
		Type:       1,
		Date:       time.Now().Format("2006-01-02T15:04:05-0700"),
		Currency:   1,
		Lang:       "rus",
		TariffCode: cdekreq.TariffCode,
		From:       request.CDEKLocation{Address: cdekreq.From},
		To:         request.CDEKLocation{Address: cdekreq.To},
		Packages:   []request.CDEKPackage{pkg},
		//Services:   services,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("failed to marshal request body", sl.Err(err))
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/calculator/tariff",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		log.Error("failed to build HTTP request", sl.Err(err))
		return nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	log.Debug("sending request to CDEK API", slog.String("url", baseURL+"/calculator/tariff"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("failed to execute HTTP request", sl.Err(err))
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	log.Debug("received response from CDEK API", slog.Int("status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		bodyStr := string(b)

		log.Error("tariff calculation failed",
			slog.Int("status_code", resp.StatusCode),
			slog.String("status", resp.Status),
			slog.String("response_body", bodyStr),
		)
		return nil, fmt.Errorf("tariff calc failed: status=%s body=%s", resp.Status, bodyStr)
	}

	var tariffs response.CDEKTariffCalcResponse
	if err := json.NewDecoder(resp.Body).Decode(&tariffs); err != nil {
		log.Error("failed to decode response", sl.Err(err))
		return nil, fmt.Errorf("decode response: %w", err)
	}

	log.Info("successfully retrieved CDEK tariffs")

	return &tariffs, nil
}

func (s *UserService) GetUserByID(userID string) (*response.UserResponse, error) {
	const op = "AuthService.GetUserByID"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("userID", userID),
	)

	log.Info("attempting to get user by id")

	model, err := s.repo.GetUserByID(userID)

	if err != nil {
		log.Error("failed to find user", sl.Err(err))
		return nil, err
	}

	resp := &response.UserResponse{
		ID:        model.ID,
		Username:  model.Username,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	log.Info("successfully retrieved user by id")

	return resp, nil
}

func (s *UserService) DeleteUser(userID string) error {
	const op = "AuthService.DeleteUser"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("userID", userID),
	)

	log.Info("attempting to delete user")

	if err := s.repo.DeleteUser(userID); err != nil {
		log.Error("failed to delete user", sl.Err(err))
		return fmt.Errorf("error deleting user from repository: %w", err)
	}
	return nil
}
