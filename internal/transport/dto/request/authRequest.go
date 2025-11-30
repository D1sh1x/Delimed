package request

// SignUpInput Структура запроса для регистрации пользователя
type SignUpInput struct {
	Username        string `json:"username" binding:"required" example:"user123"`        // Имя пользователя
	Password        string `json:"password" binding:"required,min=8" example:"password123"` // Пароль (минимум 8 символов)
	PasswordConfirm string `json:"passwordConfirm" binding:"required" example:"password123"` // Подтверждение пароля
}

// SignInInput Структура запроса для входа пользователя
type SignInInput struct {
	Username string `json:"username" binding:"required" example:"user123"` // Имя пользователя
	Password string `json:"password"  binding:"required" example:"password123"` // Пароль
}
