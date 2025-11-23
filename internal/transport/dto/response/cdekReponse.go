package response

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type CDEKServicePrice struct {
	Code            string  `json:"code"`
	Sum             float64 `json:"sum"`
	TotalSum        float64 `json:"total_sum"`
	DiscountPercent float64 `json:"discount_percent"`
	DiscountSum     float64 `json:"discount_sum"`
	VatRate         float64 `json:"vat_rate"`
	VatSum          float64 `json:"vat_sum"`
}

type CDEKErrorMessage struct {
	Code           string `json:"code"`
	AdditionalCode string `json:"additional_code"`
	Message        string `json:"message"`
}

type CDEKWarningMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CDEKDeliveryDateRange struct {
	Min string `json:"min"` // "2022-02-02"
	Max string `json:"max"` // "2022-02-04"
}

type CDEKTariffCalcResponse struct {
	DeliverySum       float64                `json:"delivery_sum"`
	PeriodMin         int                    `json:"period_min"`
	PeriodMax         int                    `json:"period_max"`
	CalendarMin       int                    `json:"calendar_min"`
	CalendarMax       int                    `json:"calendar_max"`
	WeightCalc        int                    `json:"weight_calc"`
	Services          []CDEKServicePrice     `json:"services"`
	TotalSum          float64                `json:"total_sum"`
	Currency          string                 `json:"currency"`
	Errors            []CDEKErrorMessage     `json:"errors"`
	Warnings          []CDEKWarningMessage   `json:"warnings"`
	DeliveryDateRange *CDEKDeliveryDateRange `json:"delivery_date_range"`
}

type CDEKTariffListResponse struct {
	TariffCodes []CDEKTariffListItem `json:"tariff_codes"`
	Errors      []CDEKErrorMessage   `json:"errors"`
	Warnings    []CDEKWarningMessage `json:"warnings"`
}

type CDEKTariffListItem struct {
	TariffCode        int                   `json:"tariff_code"`
	TariffName        string                `json:"tariff_name"`
	TariffDescription string                `json:"tariff_description"`
	DeliveryMode      int                   `json:"delivery_mode"`
	DeliverySum       float64               `json:"delivery_sum"`
	PeriodMin         int                   `json:"period_min"`
	PeriodMax         int                   `json:"period_max"`
	CalendarMin       int                   `json:"calendar_min"`
	CalendarMax       int                   `json:"calendar_max"`
	DeliveryDateRange CDEKDeliveryDateRange `json:"delivery_date_range"`
}
