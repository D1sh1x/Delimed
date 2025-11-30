package request

// DellinCalculatorRequest - запрос к API Деловых Линий
type DellinCalculatorRequest struct {
	AppKey   string              `json:"appkey"`
	Delivery DellinDeliveryBlock `json:"delivery"`
	Cargo    DellinCargo         `json:"cargo"`
}

// DellinDeliveryBlock - блок доставки
type DellinDeliveryBlock struct {
	DeliveryType DellinDeliveryType `json:"deliveryType"`
	Derival      DellinDerival      `json:"derival"`
	Arrival      DellinArrival      `json:"arrival"`
}

// DellinDeliveryType - тип доставки
type DellinDeliveryType struct {
	Type string `json:"type"` // "auto" / "express" / "avia" / "small" / "letter"
}

// DellinDerival - отправление
type DellinDerival struct {
	ProduceDate string              `json:"produceDate"`       // YYYY-MM-DD (обязателен)
	Variant     string              `json:"variant"`           // "address" / "terminal"
	Address     *DellinAddressBlock `json:"address,omitempty"` // если variant="address"
	Time        *DellinTimeBlock    `json:"time,omitempty"`    // если variant="address"
}

// DellinArrival - доставка
type DellinArrival struct {
	Variant string              `json:"variant"`           // "address" / "terminal"
	Address *DellinAddressBlock `json:"address,omitempty"` // если variant="address"
	Time    *DellinTimeBlock    `json:"time,omitempty"`    // если variant="address"
}

// DellinAddressBlock - адрес
type DellinAddressBlock struct {
	Search string `json:"search"` // строка адреса
}

// DellinTimeBlock - временной блок
type DellinTimeBlock struct {
	WorktimeStart string `json:"worktimeStart"` // "09:00"
	WorktimeEnd   string `json:"worktimeEnd"`   // "18:00"
}

// DellinCargo - параметры груза
type DellinCargo struct {
	Quantity    int              `json:"quantity"`            // количество мест
	Length      float64          `json:"length"`              // м
	Width       float64          `json:"width"`               // м
	Height      float64          `json:"height"`              // м
	TotalVolume float64          `json:"totalVolume"`         // м³
	TotalWeight float64          `json:"totalWeight"`         // кг
	Weight      float64          `json:"weight,omitempty"`    // вес самого тяжёлого места (если quantity>1)
	HazardClass float64          `json:"hazardClass"`         // 0 для обычного груза
	Insurance   *DellinInsurance `json:"insurance,omitempty"` // если есть страховка
}

// DellinInsurance - страхование
type DellinInsurance struct {
	StatedValue float64 `json:"statedValue"` // объявленная стоимость
	Term        bool    `json:"term"`        // страховка срока (можно false на старте)
}
