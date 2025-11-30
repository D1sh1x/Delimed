package response

import (
	"time"

	"github.com/google/uuid"
)

// UserResponse Структура ответа с информацией о пользователе
type UserResponse struct {
	ID        uuid.UUID `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"` // UUID пользователя
	Username  string    `json:"username,omitempty" example:"user123"`                        // Имя пользователя
	Role      string    `json:"role,omitempty" example:"user"`                              // Роль пользователя
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:00:00Z"`                   // Дата создания
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:00:00Z"`                   // Дата обновления
}
