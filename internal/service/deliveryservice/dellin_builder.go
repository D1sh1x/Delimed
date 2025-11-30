package deliveryservice

import (
	"delimed/internal/transport/dto/request"
	"time"
)

// BuildDellinCalcRequest - строит запрос к API Dellin из унифицированного DTO
func BuildDellinCalcRequest(appKey string, in request.DeliveryCalcRequest) request.DellinCalculatorRequest {
	// 1. Переводим размеры из см в метры
	lengthM := float64(in.LengthCm) / 100.0
	widthM := float64(in.WidthCm) / 100.0
	heightM := float64(in.HeightCm) / 100.0
	volumePerPlace := lengthM * widthM * heightM
	quantity := 1 // по умолчанию 1 место
	totalVolume := volumePerPlace * float64(quantity)
	totalWeight := in.WeightKg * float64(quantity)

	// 2. Выбираем DeliveryType (Dellin)
	deliveryTypeStr := "auto" // по умолчанию
	if in.Speed == "express" {
		deliveryTypeStr = "express"
	}

	// 3. Определяем вариант доставки (address/terminal)
	variant := "terminal"
	if in.DeliveryType == "door" {
		variant = "address"
	}

	// 4. Формируем ProduceDate
	produceDate := in.ShipmentDate
	if produceDate == "" {
		produceDate = time.Now().Format("2006-01-02")
	}

	// 5. Формируем derival
	derival := request.DellinDerival{
		ProduceDate: produceDate,
		Variant:     variant,
	}

	if variant == "address" {
		derival.Address = &request.DellinAddressBlock{
			Search: in.FromAddress,
		}
		derival.Time = &request.DellinTimeBlock{
			WorktimeStart: "09:00",
			WorktimeEnd:   "18:00",
		}
	}

	// 6. Формируем arrival
	arrival := request.DellinArrival{
		Variant: variant,
	}

	if variant == "address" {
		arrival.Address = &request.DellinAddressBlock{
			Search: in.ToAddress,
		}
		arrival.Time = &request.DellinTimeBlock{
			WorktimeStart: "09:00",
			WorktimeEnd:   "18:00",
		}
	}

	// 7. Формируем cargo
	cargo := request.DellinCargo{
		Quantity:    quantity,
		Length:      lengthM,
		Width:       widthM,
		Height:      heightM,
		TotalVolume: totalVolume,
		TotalWeight: totalWeight,
		Weight:      in.WeightKg, // самое тяжёлое место = одно место
		HazardClass: 0,
		Insurance:   buildDellinInsurance(in.ExtraServices),
	}

	// 8. Собираем запрос
	return request.DellinCalculatorRequest{
		AppKey: appKey,
		Delivery: request.DellinDeliveryBlock{
			DeliveryType: request.DellinDeliveryType{
				Type: deliveryTypeStr,
			},
			Derival: derival,
			Arrival: arrival,
		},
		Cargo: cargo,
	}
}

// buildDellinInsurance - строит блок страхования для Dellin
func buildDellinInsurance(extra request.ExtraServices) *request.DellinInsurance {
	if extra.InsuranceValue <= 0 {
		return nil
	}
	return &request.DellinInsurance{
		StatedValue: float64(extra.InsuranceValue),
		Term:        false, // на старте без страхования срока
	}
}

