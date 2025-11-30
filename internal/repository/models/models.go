package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type Order struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`

	UserID uuid.UUID `gorm:"type:uuid;not null;index"`
	User   User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Status string `gorm:"type:varchar(50);not null;index"`

	TotalAmount    int64  `gorm:"not null"`              // сумма заказа (товары + доставка), в копейках
	DeliveryAmount int64  `gorm:"not null"`              // стоимость доставки
	Currency       string `gorm:"type:char(3);not null"` // "RUB"

	DeliveryService    string `gorm:"type:varchar(50);not null"` // "cdek"
	DeliveryTariff     string `gorm:"type:varchar(50)"`
	TrackingNumber     string `gorm:"type:varchar(100);index"`
	DeliveryAddress    string `gorm:"type:text;not null"`
	DeliveryCityCode   int    `gorm:"index"`
	DeliveryPostalCode string `gorm:"type:varchar(20)"`

	Items []OrderItem `gorm:"foreignKey:OrderID"` // связи с позициями

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

type OrderItem struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key"`

	// Связь с заказом
	OrderID uuid.UUID `gorm:"type:uuid;not null;index"`
	Order   Order     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// Информация о товаре на момент покупки
	ProductID   uuid.UUID `gorm:"type:uuid;index"`         // id товара в каталоге (если есть)
	ProductSKU  string    `gorm:"type:varchar(100);index"` // артикул
	ProductName string    `gorm:"type:varchar(255);not null"`

	// Денежные поля — в минимальных единицах (копейки)
	Price       int64 `gorm:"not null"` // цена за единицу
	Quantity    int   `gorm:"not null"` // количество
	TotalAmount int64 `gorm:"not null"` // Price * Quantity

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
