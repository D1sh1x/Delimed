package domain

import "time"

// DeliveryType - тип доставки
type DeliveryType string

const (
	DeliveryTypePickup DeliveryType = "pickup" // самовывоз
	DeliveryTypeDoor   DeliveryType = "door"   // до двери
)

// DeliveryOption - единый вариант доставки от любого провайдера
// @Description Вариант доставки от провайдера (СДЭК или Деловые Линии)
type DeliveryOption struct {
	Provider     string       `json:"provider" example:"cdek"`                           // "cdek" / "dellin"
	TariffCode   string       `json:"tariff_code" example:"139"`                         // код тарифа провайдера
	Name         string       `json:"name" example:"Экспресс-лайт"`                      // человекочитаемое имя тарифа
	DeliveryType DeliveryType `json:"delivery_type" example:"door"`                      // pickup / door (deprecated)
	FromType     string        `json:"from_type,omitempty" example:"склад"`              // "склад" или "дверь" - откуда
	ToType       string        `json:"to_type,omitempty" example:"дверь"`                // "склад" или "дверь" - куда
	DeliveryMode int           `json:"delivery_mode,omitempty" example:"2"`              // delivery_mode СДЭК: 1-склад-склад, 2-склад-дверь, 3-дверь-склад, 4-дверь-дверь
	Price        int64         `json:"price" example:"150000"`                           // итоговая цена (минимальные единицы, копейки)
	Currency     string        `json:"currency" example:"RUB"`                           // "RUB"
	ETAFrom      *time.Time    `json:"eta_from,omitempty" example:"2024-01-20T00:00:00Z"` // дата "от"
	ETATo        *time.Time    `json:"eta_to,omitempty" example:"2024-01-22T00:00:00Z"`   // дата "до"
}

// FilterResult - результат фильтрации вариантов доставки
// @Description Результат расчета вариантов доставки с фильтрацией по типу
type FilterResult struct {
	Status  string           `json:"status" example:"ok"` // "ok" или "error"
	Options []DeliveryOption `json:"options"`             // отфильтрованные или fallback
}
