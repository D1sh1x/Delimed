package request

type CDEKTariffCalcRequest struct {
	Type       int           `json:"type"`               // обычно 1 — доставка
	Date       string        `json:"date,omitempty"`     // ISO-строка, можно time.Now().Format(...)
	Currency   int           `json:"currency,omitempty"` // 1 = RUB
	Lang       string        `json:"lang,omitempty"`     // "rus"
	TariffCode int           `json:"tariff_code"`        // конкретный тариф (зависит от типа и скорости)
	From       CDEKLocation  `json:"from_location"`      // откуда
	To         CDEKLocation  `json:"to_location"`        // куда
	Packages   []CDEKPackage `json:"packages"`           // габариты + вес
	Services   []CDEKService `json:"services,omitempty"` // доп. услуги
}

type CDEKTariffCalcRequestList struct {
	Type     int           `json:"type"`               // обычно 1 — доставка
	Date     string        `json:"date,omitempty"`     // ISO-строка, можно time.Now().Format(...)
	Currency int           `json:"currency,omitempty"` // 1 = RUB
	Lang     string        `json:"lang,omitempty"`     // "rus"
	From     CDEKLocation  `json:"from_location"`      // откуда
	To       CDEKLocation  `json:"to_location"`        // куда
	Packages []CDEKPackage `json:"packages"`           // габариты + вес
	Services []CDEKService `json:"services,omitempty"` // доп. услуги
}

type CDEKLocation struct {
	Code    int    `json:"code,omitempty"`    // код города СДЭК
	City    string `json:"city,omitempty"`    // можно по адресу считать, если без code
	Address string `json:"address,omitempty"` // полный адрес (если без code)
	// при желании можно добавить postcode, country_code и т.д.
}

type CDEKPackage struct {
	Weight int `json:"weight"` // граммы
	Length int `json:"length"` // см
	Width  int `json:"width"`  // см
	Height int `json:"height"` // см
}

type CDEKService struct {
	Code      string `json:"code"`                // код доп. услуги (из Приложения 3 СДЭК)
	Parameter string `json:"parameter,omitempty"` // значение параметра (зависит от услуги)
}

// CDEKRequest Структура запроса для расчета конкретного тарифа СДЭК
type CDEKRequest struct {
	TariffCode int    `json:"tariff_code" example:"139" binding:"required"` // Код тарифа СДЭК
	Weight     int    `json:"weight" example:"2500" binding:"required"`     // Вес в граммах
	Length     int    `json:"length" example:"30" binding:"required"`      // Длина в см
	Width      int    `json:"width" example:"20" binding:"required"`       // Ширина в см
	Height     int    `json:"height" example:"15" binding:"required"`     // Высота в см
	From       string `json:"from_address" example:"Москва" binding:"required"` // Адрес отправления
	To         string `json:"to_address" example:"Санкт-Петербург" binding:"required"` // Адрес доставки
}

// CDEKRequestList Структура запроса для получения списка тарифов СДЭК
type CDEKRequestList struct {
	Weight int    `json:"weight" example:"2500" binding:"required"`     // Вес в граммах
	Length int    `json:"length" example:"30" binding:"required"`       // Длина в см
	Width  int    `json:"width" example:"20" binding:"required"`        // Ширина в см
	Height int    `json:"height" example:"15" binding:"required"`       // Высота в см
	From   string `json:"from_address" example:"Москва" binding:"required"` // Адрес отправления
	To     string `json:"to_address" example:"Санкт-Петербург" binding:"required"` // Адрес доставки
}
