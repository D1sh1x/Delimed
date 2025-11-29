package deliveryservice

import (
	"delimed/internal/transport/dto/request"
	"strconv"
)

// Константы кодов услуг СДЭК
const (
	CDEKServiceInsuranceCode      = "INSURANCE"
	CDEKServiceGetUpFloorByHandCode = "GET_UP_FLOOR_BY_HAND" // Подъём на этаж (по лестнице)
	CDEKServicePackingCode         = "PACKAGE_1"             // TODO: уточнить код
	CDEKServiceDocumentsCode       = "DOCUMENTS"             // TODO: уточнить код
	CDEKServiceStorageCode         = "STORAGE"               // TODO: уточнить код
)

// MapExtraServicesToCDEK - преобразует общие дополнительные услуги в формат СДЭК
func MapExtraServicesToCDEK(services request.ExtraServices) []request.CDEKService {
	cdekServices := make([]request.CDEKService, 0)

	// Страхование
	if services.InsuranceValue > 0 {
		cdekServices = append(cdekServices, request.CDEKService{
			Code:      CDEKServiceInsuranceCode,
			Parameter: strconv.FormatInt(services.InsuranceValue, 10),
		})
	}

	// Подъём на этаж (по лестнице)
	if services.NeedUnloading {
		floors := 1 // заглушка по умолчанию
		if services.Floor != nil && *services.Floor > 0 {
			floors = *services.Floor
		}
		cdekServices = append(cdekServices, request.CDEKService{
			Code:      CDEKServiceGetUpFloorByHandCode,
			Parameter: strconv.Itoa(floors), // количество этажей
		})
	}

	// Упаковка
	if services.NeedPacking {
		cdekServices = append(cdekServices, request.CDEKService{
			Code:      CDEKServicePackingCode,
			Parameter: "", // TODO: уточнить параметры
		})
	}

	// Работа с документами
	if services.NeedDocuments {
		cdekServices = append(cdekServices, request.CDEKService{
			Code:      CDEKServiceDocumentsCode,
			Parameter: "", // TODO: уточнить параметры
		})
	}

	// Хранение
	if services.NeedStorage {
		cdekServices = append(cdekServices, request.CDEKService{
			Code:      CDEKServiceStorageCode,
			Parameter: "", // TODO: уточнить параметры
		})
	}

	return cdekServices
}

// MapExtraServicesToDellin - преобразует общие дополнительные услуги в формат Деловых Линий
// DEPRECATED: Используйте buildDellinInsurance в dellin_builder.go для страхования.
// Эта функция оставлена для будущего расширения (упаковка, документы, хранение).
func MapExtraServicesToDellin(cargo *request.DellinCargo, services request.ExtraServices) {
	// Страхование обрабатывается в buildDellinInsurance в dellin_builder.go

	// TODO: Реализовать упаковку, документы, хранение когда будут готовы структуры в API
	// Упаковка
	// if services.NeedPacking {
	//     cargo.Packing = ...
	// }

	// Работа с документами
	// if services.NeedDocuments {
	//     cargo.Documents = ...
	// }

	// Хранение
	// if services.NeedStorage {
	//     cargo.Storage = ...
	// }
}
