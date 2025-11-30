package request

// DeliveryCalcRequest - единый запрос на расчет доставки
// @Description Единый запрос для расчета вариантов доставки от всех провайдеров
type DeliveryCalcRequest struct {
	// Габариты
	LengthCm int `json:"length_cm" example:"30" binding:"required"` // Длина в см
	WidthCm  int `json:"width_cm" example:"20" binding:"required"`  // Ширина в см
	HeightCm int `json:"height_cm" example:"15" binding:"required"` // Высота в см

	// Вес
	WeightKg float64 `json:"weight_kg" example:"2.5" binding:"required"` // Вес в кг

	// Тип доставки: самовывоз / дверь
	// pickup — самовывоз из ПВЗ/терминала
	// door   — доставка до двери
	DeliveryType string `json:"delivery_type" example:"door" binding:"required"` // "pickup" или "door"

	// Тариф по скорости (сейчас нужна только фильтрация по типу доставки,
	// но speed лучше сразу оставить в модели на будущее)
	Speed string `json:"speed" example:"economy"` // "economy" / "express" / "urgent"

	// Адреса
	FromAddress string `json:"from_address" example:"Москва, ул. Ленина, д. 1" binding:"required"`         // Адрес отправления
	ToAddress   string `json:"to_address" example:"Санкт-Петербург, Невский пр., д. 1" binding:"required"` // Адрес доставки

	// Дата отгрузки (если нет — можно подставлять today / today+1 на бэкенде)
	ShipmentDate string `json:"shipment_date" example:"2024-01-15"` // "YYYY-MM-DD"

	// Доп услуги
	ExtraServices ExtraServices `json:"extra_services"` // Дополнительные услуги
}

// ExtraServices - дополнительные услуги
// @Description Дополнительные услуги для расчета доставки
type ExtraServices struct {
	InsuranceValue int64 `json:"insurance_value" example:"5000"`         // объявленная стоимость, 0 если не нужно
	NeedPacking    bool  `json:"need_packing" example:"false"`           // упаковка
	NeedCourier    bool  `json:"need_courier" example:"false"`           // курьеры/адресная доставка/грузчики
	NeedDocuments  bool  `json:"need_documents" example:"false"`         // работа с документами
	NeedStorage    bool  `json:"need_storage" example:"false"`           // хранение
	NeedUnloading  bool  `json:"need_unloading" example:"false"`         // разгрузочные работы / подъём на этаж
	Floor          *int  `json:"floor,omitempty" example:"5"`            // этаж доставки (если не указан, используется 1)
	HasElevator    *bool `json:"has_elevator,omitempty" example:"false"` // наличие лифта (для будущего использования)
}
