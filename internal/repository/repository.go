package repository

import "delimed/internal/repository/models"

type RepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)

	GetUserByID(userID string) (*models.User, error)
	UpdateUser(userID string, updates models.User) error
	DeleteUser(userID string) error

	// Order CRUD
	CreateOrder(order *models.Order) error
	GetOrderByID(orderID string) (*models.Order, error)
	GetOrderWithItems(orderID string) (*models.Order, error)
	GetOrdersByUserID(userID string) ([]*models.Order, error)
	UpdateOrder(orderID string, updates models.Order) error
	DeleteOrder(orderID string) error

	// OrderItem CRUD
	CreateOrderItem(item *models.OrderItem) error
	GetOrderItems(orderID string) ([]*models.OrderItem, error)
	GetOrderItemByID(itemID string) (*models.OrderItem, error)
	UpdateOrderItem(itemID string, updates models.OrderItem) error
	DeleteOrderItem(itemID string) error
}
