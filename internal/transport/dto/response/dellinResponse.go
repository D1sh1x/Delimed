package response

// DellinCalculatorResponse - ответ от API Деловых Линий
type DellinCalculatorResponse struct {
	Metadata DellinMetadata       `json:"metadata"`
	Data     DellinCalculatorData `json:"data"`
}

// DellinMetadata - метаданные ответа
type DellinMetadata struct {
	Status      int    `json:"status"`
	GeneratedAt string `json:"generated_at"`
}

// DellinCalculatorData - данные расчета
type DellinCalculatorData struct {
	Price                  float64            `json:"price"`
	PriceMinimal           string             `json:"priceMinimal"`
	AvailableDeliveryTypes map[string]float64 `json:"availableDeliveryTypes"` // "auto": 480.0, "express": 620.0...
	OrderDates             DellinOrderDates   `json:"orderDates"`
}

// DellinOrderDates - даты доставки
type DellinOrderDates struct {
	Pickup                    *string `json:"pickup"`
	ArrivalToOspReceiver      *string `json:"arrivalToOspReceiver"`
	ArrivalToAirport          *string `json:"arrivalToAirport"`
	ArrivalToAirportMax       *string `json:"arrivalToAirportMax"`
	GiveoutFromOspReceiver    *string `json:"giveoutFromOspReceiver"`
	GiveoutFromOspReceiverMax *string `json:"giveoutFromOspReceiverMax"`
}

// DellinError - ошибка API (для обработки ошибок)
type DellinError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
