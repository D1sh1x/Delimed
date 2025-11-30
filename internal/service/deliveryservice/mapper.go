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
		// Определяем типы доставки по delivery_mode
		// delivery_mode: 1 - склад-склад, 2 - склад-дверь, 3 - дверь-склад, 4 - дверь-дверь
		deliveryType, arrivalType := determineCDEKDeliveryTypes(tariff.DeliveryMode, tariff.TariffName)

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
			ArrivalType:  arrivalType,
			Price:        price,
			Currency:     "RUB",
			ETAFrom:      etaFrom,
			ETATo:        etaTo,
		}

		options = append(options, option)
	}

	return options
}

// determineCDEKDeliveryTypes - определяет типы доставки для СДЭК (откуда и куда)
// delivery_mode: 1 - склад-склад, 2 - склад-дверь, 3 - дверь-склад, 4 - дверь-дверь
// Возвращает (deliveryType, arrivalType) - откуда и куда
func determineCDEKDeliveryTypes(deliveryMode int, tariffName string) (domain.DeliveryType, domain.DeliveryType) {
	// Если есть delivery_mode, используем его
	switch deliveryMode {
	case 1:
		// склад-склад
		return domain.DeliveryTypePickup, domain.DeliveryTypePickup
	case 2:
		// склад-дверь
		return domain.DeliveryTypePickup, domain.DeliveryTypeDoor
	case 3:
		// дверь-склад
		return domain.DeliveryTypeDoor, domain.DeliveryTypePickup
	case 4:
		// дверь-дверь
		return domain.DeliveryTypeDoor, domain.DeliveryTypeDoor
	}

	// Fallback: пытаемся определить по названию тарифа
	nameLower := strings.ToLower(tariffName)
	hasSklad := strings.Contains(nameLower, "склад")
	hasDoor := strings.Contains(nameLower, "дверь")

	if hasSklad && !hasDoor {
		// только склад
		return domain.DeliveryTypePickup, domain.DeliveryTypePickup
	}
	if hasDoor && !hasSklad {
		// только дверь
		return domain.DeliveryTypeDoor, domain.DeliveryTypeDoor
	}
	if hasSklad && hasDoor {
		// оба типа - пытаемся определить порядок
		// По умолчанию считаем склад-дверь
		return domain.DeliveryTypePickup, domain.DeliveryTypeDoor
	}

	// По умолчанию считаем door-door, если не указано иное
	return domain.DeliveryTypeDoor, domain.DeliveryTypeDoor
}

// MapDellinCalcToOptions - преобразует ответ Деловых Линий в единый формат DeliveryOption
// arrivalVariant: "address" (дверь) или "terminal" (склад) - куда доставляется
func MapDellinCalcToOptions(resp response.DellinCalculatorResponse, arrivalVariant string) []domain.DeliveryOption {
	options := make([]domain.DeliveryOption, 0)

	// Проверяем статус ответа
	// Status может быть 0 (успех по умолчанию) или 200
	// Если статус явно указан и не равен 0/200, возвращаем пустой список
	if resp.Metadata.Status != 0 && resp.Metadata.Status != 200 {
		return options
	}

	// Определяем arrivalType на основе variant
	var arrivalType domain.DeliveryType
	if arrivalVariant == "address" {
		arrivalType = domain.DeliveryTypeDoor
	} else {
		arrivalType = domain.DeliveryTypePickup // terminal
	}

	// Если есть availableDeliveryTypes, создаем варианты для каждого типа
	if len(resp.Data.AvailableDeliveryTypes) > 0 {
		for deliveryTypeStr, price := range resp.Data.AvailableDeliveryTypes {
			// Определяем тип доставки (по умолчанию door, так как считаем адрес→адрес)
			// Но на самом деле это зависит от variant в запросе
			deliveryType := domain.DeliveryTypeDoor
			// Если variant был "terminal", то и отправление тоже terminal
			if arrivalVariant == "terminal" {
				deliveryType = domain.DeliveryTypePickup
			}

			// Название тарифа
			tariffName := getDellinTariffName(deliveryTypeStr)

			option := domain.DeliveryOption{
				Provider:     "dellin",
				TariffCode:   deliveryTypeStr,
				Name:         tariffName,
				DeliveryType: deliveryType,
				ArrivalType:  arrivalType,
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
		// Если variant был "terminal", то и отправление тоже terminal
		if arrivalVariant == "terminal" {
			deliveryType = domain.DeliveryTypePickup
		}

		tariffCode := resp.Data.PriceMinimal
		if tariffCode == "" {
			tariffCode = "auto"
		}

		option := domain.DeliveryOption{
			Provider:     "dellin",
			TariffCode:   tariffCode,
			Name:         fmt.Sprintf("Деловые линии (%s)", tariffCode),
			DeliveryType: deliveryType,
			ArrivalType:  arrivalType,
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

// FilterOptionsByDeliveryType - фильтрует варианты доставки по типу
// Фильтрует по второму типу доставки (arrival type - куда доставляется):
// - если requestedType = "door" - оставляем только тарифы где второй тип "door" (склад-дверь, дверь-дверь)
// - если requestedType = "pickup" - оставляем только тарифы где второй тип "pickup" (склад-склад, дверь-склад)
func FilterOptionsByDeliveryType(all []domain.DeliveryOption, requestedType domain.DeliveryType) domain.FilterResult {
	filtered := make([]domain.DeliveryOption, 0)

	for _, option := range all {
		nameLower := strings.ToLower(option.Name)
		if strings.Contains(nameLower, "постомат") {
			continue
		}

		// Фильтруем по второму типу доставки (arrival type - куда доставляется)
		// Если ArrivalType не установлен, используем старую логику (DeliveryType) для обратной совместимости
		arrivalType := option.ArrivalType
		if arrivalType == "" {
			// Fallback: если ArrivalType не установлен, используем DeliveryType
			arrivalType = option.DeliveryType
		}

		// Проверяем, что второй тип доставки (куда) совпадает с запрошенным
		if arrivalType == requestedType {
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
