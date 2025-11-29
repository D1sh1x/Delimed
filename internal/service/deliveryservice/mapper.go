package deliveryservice

import (
	"delimed/internal/domain"
	"delimed/internal/transport/dto/response"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// MapCDEKTarifflistToOptions - преобразует ответ СДЭК в единый формат DeliveryOption
func MapCDEKTarifflistToOptions(resp response.CDEKTariffListResponse) []domain.DeliveryOption {
	options := make([]domain.DeliveryOption, 0, len(resp.TariffCodes))

	for _, tariff := range resp.TariffCodes {
		// Определяем тип доставки по delivery_mode
		// delivery_mode: 1 - склад-склад, 2 - склад-дверь, 3 - дверь-склад, 4 - дверь-дверь
		deliveryType := determineCDEKDeliveryType(tariff.DeliveryMode, tariff.TariffName)

		// Определяем from_type и to_type на основе delivery_mode
		fromType, toType := getCDEKFromToTypes(tariff.DeliveryMode)

		// Преобразуем цену из рублей в копейки
		price := int64(tariff.DeliverySum * 100)

		// Парсим даты доставки
		var etaFrom, etaTo *time.Time
		if tariff.DeliveryDateRange.Min != "" {
			if t, err := time.Parse("2006-01-02", tariff.DeliveryDateRange.Min); err == nil {
				etaFrom = &t
			}
		}
		if tariff.DeliveryDateRange.Max != "" {
			if t, err := time.Parse("2006-01-02", tariff.DeliveryDateRange.Max); err == nil {
				etaTo = &t
			}
		}

		option := domain.DeliveryOption{
			Provider:     "cdek",
			TariffCode:   strconv.Itoa(tariff.TariffCode),
			Name:         tariff.TariffName,
			DeliveryType: deliveryType,
			FromType:     fromType,
			ToType:       toType,
			DeliveryMode: tariff.DeliveryMode,
			Price:        price,
			Currency:     "RUB",
			ETAFrom:      etaFrom,
			ETATo:        etaTo,
		}

		options = append(options, option)
	}

	return options
}

// determineCDEKDeliveryType - определяет тип доставки для СДЭК
// delivery_mode: 1 - склад-склад (pickup), 2 - склад-дверь (door), 3 - дверь-склад (door), 4 - дверь-дверь (door)
func determineCDEKDeliveryType(deliveryMode int, tariffName string) domain.DeliveryType {
	// Если есть delivery_mode, используем его
	if deliveryMode == 1 {
		return domain.DeliveryTypePickup // склад-склад
	}
	if deliveryMode == 2 || deliveryMode == 3 || deliveryMode == 4 {
		return domain.DeliveryTypeDoor // склад-дверь, дверь-склад, дверь-дверь
	}

	// Fallback: пытаемся определить по названию тарифа
	nameLower := strings.ToLower(tariffName)
	if strings.Contains(nameLower, "склад") && !strings.Contains(nameLower, "дверь") {
		return domain.DeliveryTypePickup
	}
	if strings.Contains(nameLower, "дверь") {
		return domain.DeliveryTypeDoor
	}

	// По умолчанию считаем door, если не указано иное
	return domain.DeliveryTypeDoor
}

// getCDEKFromToTypes - определяет from_type и to_type на основе delivery_mode СДЭК
// delivery_mode: 1 - склад-склад, 2 - склад-дверь, 3 - дверь-склад, 4 - дверь-дверь
func getCDEKFromToTypes(deliveryMode int) (fromType, toType string) {
	switch deliveryMode {
	case 1:
		return "склад", "склад" // склад-склад
	case 2:
		return "склад", "дверь" // склад-дверь
	case 3:
		return "дверь", "склад" // дверь-склад
	case 4:
		return "дверь", "дверь" // дверь-дверь
	default:
		// По умолчанию
		return "склад", "дверь"
	}
}

// MapDellinCalcToOptions - преобразует ответ Деловых Линий в единый формат DeliveryOption
// to используется для определения to_type, from всегда "склад"
func MapDellinCalcToOptions(resp response.DellinCalculatorResponse, to string) []domain.DeliveryOption {
	options := make([]domain.DeliveryOption, 0)

	// Проверяем статус ответа
	// Status может быть 0 (успех по умолчанию) или 200
	// Если статус явно указан и не равен 0/200, возвращаем пустой список
	if resp.Metadata.Status != 0 && resp.Metadata.Status != 200 {
		return options
	}

	// Если есть availableDeliveryTypes, создаем варианты для каждого типа
	if len(resp.Data.AvailableDeliveryTypes) > 0 {
		for deliveryTypeStr, price := range resp.Data.AvailableDeliveryTypes {
			// Определяем тип доставки (по умолчанию door, так как считаем адрес→адрес)
			deliveryType := domain.DeliveryTypeDoor

			// Название тарифа
			tariffName := getDellinTariffName(deliveryTypeStr)

			// Для Деловых Линий from всегда "склад", to берем из параметра
			fromType, toType := "склад", to

			option := domain.DeliveryOption{
				Provider:     "dellin",
				TariffCode:   deliveryTypeStr,
				Name:         tariffName,
				DeliveryType: deliveryType,
				FromType:     fromType,
				ToType:       toType,
				Price:        int64(math.Round(price * 100)), // рубли в копейки
				Currency:     "RUB",
			}

			// Парсим даты доставки
			etaFrom := parseDellinDate(resp.Data.OrderDates.GiveoutFromOspReceiver)
			etaTo := parseDellinDate(resp.Data.OrderDates.GiveoutFromOspReceiverMax)
			if etaFrom != nil {
				option.ETAFrom = etaFrom
			}
			if etaTo != nil {
				option.ETATo = etaTo
			}

			options = append(options, option)
		}
	} else {
		// Если нет availableDeliveryTypes, используем общий price
		deliveryType := domain.DeliveryTypeDoor // по умолчанию door (адрес→адрес)

		tariffCode := resp.Data.PriceMinimal
		if tariffCode == "" {
			tariffCode = "auto"
		}

		// Для Деловых Линий from всегда "склад", to берем из параметра
		fromType, toType := "склад", to

		option := domain.DeliveryOption{
			Provider:     "dellin",
			TariffCode:   tariffCode,
			Name:         fmt.Sprintf("Деловые линии (%s)", tariffCode),
			DeliveryType: deliveryType,
			FromType:     fromType,
			ToType:       toType,
			Price:        int64(math.Round(resp.Data.Price * 100)), // рубли в копейки
			Currency:     "RUB",
		}

		// Парсим даты доставки
		etaFrom := parseDellinDate(resp.Data.OrderDates.GiveoutFromOspReceiver)
		etaTo := parseDellinDate(resp.Data.OrderDates.GiveoutFromOspReceiverMax)
		if etaFrom != nil {
			option.ETAFrom = etaFrom
		}
		if etaTo != nil {
			option.ETATo = etaTo
		}

		options = append(options, option)
	}

	return options
}

// getDellinTariffName - возвращает человекочитаемое название тарифа Dellin
func getDellinTariffName(tariffCode string) string {
	names := map[string]string{
		"auto":    "Деловые линии (Авто)",
		"express": "Деловые линии (Экспресс)",
		"avia":    "Деловые линии (Авиа)",
		"small":   "Деловые линии (Малогабарит)",
		"letter":  "Деловые линии (Письмо)",
	}
	if name, ok := names[tariffCode]; ok {
		return name
	}
	return fmt.Sprintf("Деловые линии (%s)", tariffCode)
}

// parseDellinDate - парсит дату из формата Dellin
func parseDellinDate(dateStr *string) *time.Time {
	if dateStr == nil || *dateStr == "" {
		return nil
	}

	// Пробуем разные форматы
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, *dateStr); err == nil {
			return &t
		}
	}

	return nil
}

// FilterOptionsByDeliveryType - фильтрует варианты доставки по типу (deprecated, используйте FilterOptionsByFromTo)
func FilterOptionsByDeliveryType(all []domain.DeliveryOption, requestedType domain.DeliveryType) domain.FilterResult {
	filtered := make([]domain.DeliveryOption, 0)

	for _, option := range all {
		if option.DeliveryType == requestedType {
			filtered = append(filtered, option)
		}
	}

	if len(filtered) > 0 {
		return domain.FilterResult{
			Status:  "ok",
			Options: filtered,
		}
	}

	// Если нет подходящих вариантов, возвращаем все как fallback
	return domain.FilterResult{
		Status:  "error",
		Options: all,
	}
}

// FilterOptionsByFromTo - фильтрует варианты доставки по from/to и исключает тарифы с "постомат"
func FilterOptionsByFromTo(all []domain.DeliveryOption, from, to string) domain.FilterResult {
	filtered := make([]domain.DeliveryOption, 0)

	// Нормализуем значения from и to
	from = strings.ToLower(strings.TrimSpace(from))
	to = strings.ToLower(strings.TrimSpace(to))

	for _, option := range all {
		// Исключаем тарифы с "постомат" в названии (case-insensitive)
		nameLower := strings.ToLower(option.Name)
		if strings.Contains(nameLower, "постомат") {
			continue
		}

		// Фильтруем по from и to
		optionFrom := strings.ToLower(strings.TrimSpace(option.FromType))
		optionTo := strings.ToLower(strings.TrimSpace(option.ToType))

		if optionFrom == from && optionTo == to {
			filtered = append(filtered, option)
		}
	}

	if len(filtered) > 0 {
		return domain.FilterResult{
			Status:  "ok",
			Options: filtered,
		}
	}

	// Если нет подходящих вариантов, возвращаем все как fallback
	return domain.FilterResult{
		Status:  "error",
		Options: all,
	}
}

// FilterCDEKBySpeed - фильтрует тарифы СДЭК по скорости доставки
// Применяется только к тарифам СДЭК, тарифы Деловых Линий остаются без изменений
func FilterCDEKBySpeed(options []domain.DeliveryOption, speed string) []domain.DeliveryOption {
	// Если speed не указан, возвращаем все варианты
	if speed == "" {
		return options
	}

	speed = strings.ToLower(strings.TrimSpace(speed))
	filtered := make([]domain.DeliveryOption, 0)

	// Определяем разрешенные тарифы в зависимости от speed
	var allowedTariffs map[string]bool

	switch speed {
	case "economy":
		// Экономичная доставка: 62, 121, 122, 123, 748, 749, 750, 751
		allowedTariffs = map[string]bool{
			"62":  true,
			"121": true,
			"122": true,
			"123": true,
			"748": true,
			"749": true,
			"750": true,
			"751": true,
		}
	case "express":
		// Экспресс доставка: 480, 481, 482, 483
		allowedTariffs = map[string]bool{
			"480": true,
			"481": true,
			"482": true,
			"483": true,
		}
	default:
		// Если speed не economy и не express, возвращаем все варианты
		return options
	}

	// Фильтруем варианты
	for _, option := range options {
		// Для тарифов СДЭК применяем фильтрацию по speed
		if option.Provider == "cdek" {
			// Проверяем, есть ли тариф в списке разрешенных
			if allowedTariffs[option.TariffCode] {
				filtered = append(filtered, option)
			}
		} else {
			// Для других провайдеров (Деловые Линии) оставляем без изменений
			filtered = append(filtered, option)
		}
	}

	return filtered
}
