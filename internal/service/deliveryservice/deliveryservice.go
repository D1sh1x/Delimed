package deliveryservice

import (
	"bytes"
	"context"
	"delimed/internal/domain"
	"delimed/internal/transport/dto/request"
	"delimed/internal/transport/dto/response"
	"delimed/internal/utils/cdek"
	sl "delimed/internal/utils/logger"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"time"
)

type DeliveryServiceInterface interface {
	CalculateDeliveryOptions(ctx context.Context, req request.DeliveryCalcRequest) (domain.FilterResult, error)
}

type DeliveryService struct {
	logger *slog.Logger
	// Конфигурация API
	cdekClientID     string
	cdekClientSecret string
	dellinAppKey     string
}

func NewDeliveryService(logger *slog.Logger, cdekClientID, cdekClientSecret, dellinAppKey string) DeliveryServiceInterface {
	return &DeliveryService{
		logger:           logger,
		cdekClientID:     cdekClientID,
		cdekClientSecret: cdekClientSecret,
		dellinAppKey:     dellinAppKey,
	}
}

// CalculateDeliveryOptions - единый расчет вариантов доставки от всех провайдеров
func (s *DeliveryService) CalculateDeliveryOptions(ctx context.Context, req request.DeliveryCalcRequest) (domain.FilterResult, error) {
	const op = "DeliveryService.CalculateDeliveryOptions"

	log := s.logger.With(
		slog.String("op", op),
		slog.String("delivery_type", req.DeliveryType),
	)

	log.Info("calculating delivery options from all providers")

	// Собираем варианты от всех провайдеров параллельно
	allOptions := make([]domain.DeliveryOption, 0)

	// 1. Получаем варианты от СДЭК
	cdekOptions, err := s.getCDEKOptions(ctx, req)
	if err != nil {
		log.Warn("failed to get CDEK options", sl.Err(err))
		// Продолжаем работу даже если один провайдер не ответил
	} else {
		allOptions = append(allOptions, cdekOptions...)
	}

	// 2. Получаем варианты от Деловых Линий
	dellinOptions, err := s.getDellinOptions(ctx, req)
	if err != nil {
		log.Warn("failed to get Dellin options", sl.Err(err))
		// Продолжаем работу даже если один провайдер не ответил
	} else {
		allOptions = append(allOptions, dellinOptions...)
	}

	// 3. Фильтруем по типу доставки
	requestedType := domain.DeliveryType(req.DeliveryType)
	if requestedType != domain.DeliveryTypePickup && requestedType != domain.DeliveryTypeDoor {
		// Если тип не указан или неверный, используем door по умолчанию
		requestedType = domain.DeliveryTypeDoor
	}

	result := FilterOptionsByDeliveryType(allOptions, requestedType)

	// 4. Фильтруем тарифы СДЭК по speed (последняя фильтрация)
	if req.Speed != "" {
		result.Options = FilterCDEKBySpeed(result.Options, req.Speed)
	}

	sort.Slice(result.Options, func(i, j int) bool {
		return result.Options[i].Price < result.Options[j].Price
	})

	log.Info("delivery options calculated",
		slog.Int("total_options", len(allOptions)),
		slog.Int("filtered_options", len(result.Options)),
		slog.String("status", result.Status),
	)

	return result, nil
}

// getCDEKOptions - получает варианты доставки от СДЭК
func (s *DeliveryService) getCDEKOptions(ctx context.Context, req request.DeliveryCalcRequest) ([]domain.DeliveryOption, error) {
	const op = "DeliveryService.getCDEKOptions"

	log := s.logger.With(slog.String("op", op))

	// Получаем токен
	token, err := cdek.GetCDEKToken(ctx, s.cdekClientID, s.cdekClientSecret)
	if err != nil {
		log.Error("failed to get CDEK token", sl.Err(err))
		return nil, fmt.Errorf("get token: %w", err)
	}

	// Преобразуем вес из кг в граммы
	weightGrams := int(req.WeightKg * 1000)

	// Формируем пакет
	pkg := request.CDEKPackage{
		Weight: weightGrams,
		Length: req.LengthCm,
		Width:  req.WidthCm,
		Height: req.HeightCm,
	}

	// Преобразуем дополнительные услуги
	services := MapExtraServicesToCDEK(req.ExtraServices)

	// Формируем запрос
	reqBody := request.CDEKTariffCalcRequestList{
		Type:     1,
		Date:     time.Now().Format("2006-01-02T15:04:05-0700"),
		Currency: 1,
		Lang:     "rus",
		From:     request.CDEKLocation{Address: req.FromAddress},
		To:       request.CDEKLocation{Address: req.ToAddress},
		Packages: []request.CDEKPackage{pkg},
		Services: services,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Error("failed to marshal request body", sl.Err(err))
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	// Выполняем запрос
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.edu.cdek.ru/v2/calculator/tarifflist",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		log.Error("failed to build HTTP request", sl.Err(err))
		return nil, fmt.Errorf("build request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Error("failed to execute HTTP request", sl.Err(err))
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		bodyStr := string(b)
		log.Error("CDEK API returned error",
			slog.Int("status_code", resp.StatusCode),
			slog.String("response_body", bodyStr),
		)
		return nil, fmt.Errorf("CDEK API error: status=%s body=%s", resp.Status, bodyStr)
	}

	var cdekResp response.CDEKTariffListResponse
	if err := json.NewDecoder(resp.Body).Decode(&cdekResp); err != nil {
		log.Error("failed to decode CDEK response", sl.Err(err))
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Преобразуем в единый формат
	options := MapCDEKTarifflistToOptions(cdekResp)

	return options, nil
}

// getDellinOptions - получает варианты доставки от Деловых Линий
func (s *DeliveryService) getDellinOptions(ctx context.Context, req request.DeliveryCalcRequest) ([]domain.DeliveryOption, error) {
	const op = "DeliveryService.getDellinOptions"

	log := s.logger.With(slog.String("op", op))

	// Используем dev-ключ, если не указан в конфиге
	appKey := s.dellinAppKey
	if appKey == "" {
		appKey = "1E211F21-EE54-4EE4-B33E-F889DC83383F"
	}

	// Строим запрос к Dellin API
	dellinReq := BuildDellinCalcRequest(appKey, req)

	bodyBytes, err := json.Marshal(dellinReq)
	if err != nil {
		log.Error("failed to marshal Dellin request body", sl.Err(err))
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	// Логируем запрос для отладки
	log.Debug("Dellin request",
		slog.String("request_body", string(bodyBytes)),
	)

	// Выполняем запрос
	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.dellin.ru/v2/calculator.json",
		bytes.NewReader(bodyBytes),
	)
	if err != nil {
		log.Error("failed to build Dellin HTTP request", sl.Err(err))
		return nil, fmt.Errorf("build request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Error("failed to execute Dellin HTTP request", sl.Err(err))
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("failed to read Dellin response body", sl.Err(err))
		return nil, fmt.Errorf("read response: %w", err)
	}

	bodyStr := string(responseBody)

	if resp.StatusCode != http.StatusOK {
		log.Error("Dellin API returned error",
			slog.Int("status_code", resp.StatusCode),
			slog.String("response_body", bodyStr),
		)
		return nil, fmt.Errorf("dellin API error: status=%s body=%s", resp.Status, bodyStr)
	}

	// Логируем сырой ответ для отладки
	log.Debug("Dellin response",
		slog.String("response_body", bodyStr),
	)

	var dellinResp response.DellinCalculatorResponse
	if err := json.Unmarshal(responseBody, &dellinResp); err != nil {
		log.Error("failed to decode Dellin response",
			sl.Err(err),
			slog.String("response_body", bodyStr),
		)
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Извлекаем arrival variant из запроса для правильной фильтрации
	arrivalVariant := dellinReq.Delivery.Arrival.Variant

	// Преобразуем в единый формат
	options := MapDellinCalcToOptions(dellinResp, arrivalVariant)

	return options, nil
}
