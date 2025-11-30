package handler

import (
	"delimed/internal/transport/dto/request"
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterUserHandler godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает новую учетную запись пользователя
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.SignUpInput  true  "Данные для регистрации"
// @Success      201      {object}  map[string]string    "Регистрация успешна"
// @Failure      400      {object}  map[string]string    "Неверный формат входных данных"
// @Failure      409      {object}  map[string]string    "Пользователь уже существует"
// @Router       /register [post]
func (h *Handler) RegisterUserHandler(c echo.Context) error {

	var req request.SignUpInput

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
	}

	if err := h.service.Auth().RegisterUser(req); err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "Registration successful."})
}

// LoginUserHandler godoc
// @Summary      Вход пользователя
// @Description  Аутентификация пользователя и получение JWT токена
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      request.SignInInput  true  "Данные для входа"
// @Success      200      {object}  map[string]string    "Токен доступа"
// @Failure      400      {object}  map[string]string    "Неверный формат входных данных"
// @Failure      401      {object}  map[string]string    "Неверные учетные данные"
// @Router       /login [post]
func (h *Handler) LoginUserHandler(c echo.Context) error {
	var req request.SignInInput
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
	}

	token, err := h.service.Auth().LoginUser(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
