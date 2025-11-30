package handler

import (
	"delimed/internal/transport/dto/request"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetTariffsList godoc
// @Summary      Получить список тарифов СДЭК
// @Description  Возвращает список доступных тарифов СДЭК для указанных параметров доставки
// @Tags         delivery
// @Accept       json
// @Produce      json
// @Param        request  body      request.CDEKRequestList  true  "Параметры для расчета тарифов"
// @Success      200      {object}  response.CDEKTariffListResponse  "Список тарифов"
// @Failure      400      {object}  map[string]string                "Неверный формат входных данных"
// @Failure      500      {object}  map[string]string                "Внутренняя ошибка сервера"
// @Router       /tariffslist [get]
func (h *Handler) GetTariffsList(c echo.Context) error {
	var req request.CDEKRequestList

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input format",
		})
	}

	ctx := c.Request().Context()

	resp, err := h.service.User().GetTariffsList(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetTariffs godoc
// @Summary      Получить расчет тарифа СДЭК
// @Description  Возвращает детальный расчет стоимости доставки по конкретному тарифу СДЭК
// @Tags         delivery
// @Accept       json
// @Produce      json
// @Param        request  body      request.CDEKRequest  true  "Параметры для расчета тарифа"
// @Success      200      {object}  response.CDEKTariffCalcResponse  "Расчет тарифа"
// @Failure      400      {object}  map[string]string                 "Неверный формат входных данных"
// @Failure      500      {object}  map[string]string                 "Внутренняя ошибка сервера"
// @Router       /tariffs [get]
func (h *Handler) GetTariffs(c echo.Context) error {
	var req request.CDEKRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input format",
		})
	}

	ctx := c.Request().Context()

	resp, err := h.service.User().GetTarifs(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

// GetUserProfileHandler godoc
// @Summary      Получить профиль пользователя
// @Description  Возвращает информацию о текущем аутентифицированном пользователе
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  response.UserResponse  "Профиль пользователя"
// @Failure      401      {object}  map[string]string      "Пользователь не аутентифицирован"
// @Failure      404      {object}  map[string]string      "Пользователь не найден"
// @Failure      500      {object}  map[string]string      "Внутренняя ошибка сервера"
// @Router       /api/user [get]
func (h *Handler) GetUserProfileHandler(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	user, err := h.service.User().GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUserProfileHandler godoc
// @Summary      Удалить профиль пользователя
// @Description  Удаляет учетную запись текущего аутентифицированного пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200      {object}  map[string]string  "Пользователь успешно удален"
// @Failure      401      {object}  map[string]string  "Пользователь не аутентифицирован"
// @Failure      404      {object}  map[string]string  "Пользователь не найден"
// @Failure      500      {object}  map[string]string  "Внутренняя ошибка сервера"
// @Router       /api/user [delete]
func (h *Handler) DeleteUserProfileHandler(c echo.Context) error {
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "User not authenticated"})
	}

	userID, ok := userIDInterface.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	err := h.service.User().DeleteUser(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "User has deleted"})
}

// CalculateDeliveryOptions godoc
// @Summary      Расчет вариантов доставки
// @Description  Возвращает варианты доставки от всех провайдеров (СДЭК, Деловые Линии) с фильтрацией по типу доставки
// @Tags         delivery
// @Accept       json
// @Produce      json
// @Param        request  body      request.DeliveryCalcRequest  true  "Параметры для расчета доставки"
// @Success      200      {object}  domain.FilterResult           "Варианты доставки"
// @Failure      400      {object}  map[string]string              "Неверный формат входных данных"
// @Failure      500      {object}  map[string]string              "Внутренняя ошибка сервера"
// @Router       /delivery/calculate [post]
func (h *Handler) CalculateDeliveryOptions(c echo.Context) error {
	var req request.DeliveryCalcRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid input format",
		})
	}

	ctx := c.Request().Context()

	result, err := h.service.Delivery().CalculateDeliveryOptions(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, result)
}
