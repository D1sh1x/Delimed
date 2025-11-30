package postgres

import (
	"delimed/internal/repository"
	"delimed/internal/repository/models"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	DB *gorm.DB
}

func NewStorage(dsn string,
	connMaxIdleTime, connMaxLifetime time.Duration,
	maxOpenConns, maxIdleConns int,
) (repository.RepositoryInterface, error) {

	newLogger := logger.Default
	newLogger = newLogger.LogMode(logger.Info)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get *sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Order{}, &models.OrderItem{}); err != nil {
		return nil, fmt.Errorf("failed migrations: %w", err)
	}

	return &Storage{
		DB: db,
	}, nil
}

// User CRUD methods

func (s *Storage) CreateUser(user *models.User) error {
	if err := s.DB.Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *Storage) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return &user, nil
}

func (s *Storage) GetUserByID(userID string) (*models.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	var user models.User
	if err := s.DB.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

func (s *Storage) UpdateUser(userID string, updates models.User) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	result := s.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (s *Storage) DeleteUser(userID string) error {
	id, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	result := s.DB.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// Order CRUD methods

func (s *Storage) CreateOrder(order *models.Order) error {
	if err := s.DB.Create(order).Error; err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}
	return nil
}

func (s *Storage) GetOrderByID(orderID string) (*models.Order, error) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	var order models.Order
	if err := s.DB.First(&order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get order by ID: %w", err)
	}
	return &order, nil
}

func (s *Storage) GetOrderWithItems(orderID string) (*models.Order, error) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	var order models.Order
	if err := s.DB.Preload("Items").First(&order, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get order with items: %w", err)
	}
	return &order, nil
}

func (s *Storage) GetOrdersByUserID(userID string) ([]*models.Order, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	var orders []*models.Order
	if err := s.DB.Where("user_id = ?", id).Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to get orders by user ID: %w", err)
	}
	return orders, nil
}

func (s *Storage) UpdateOrder(orderID string, updates models.Order) error {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	result := s.DB.Model(&models.Order{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

func (s *Storage) DeleteOrder(orderID string) error {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return fmt.Errorf("invalid order ID format: %w", err)
	}

	result := s.DB.Delete(&models.Order{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete order: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}

// OrderItem CRUD methods

func (s *Storage) CreateOrderItem(item *models.OrderItem) error {
	if err := s.DB.Create(item).Error; err != nil {
		return fmt.Errorf("failed to create order item: %w", err)
	}
	return nil
}

func (s *Storage) GetOrderItems(orderID string) ([]*models.OrderItem, error) {
	id, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID format: %w", err)
	}

	var items []*models.OrderItem
	if err := s.DB.Where("order_id = ?", id).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}
	return items, nil
}

func (s *Storage) GetOrderItemByID(itemID string) (*models.OrderItem, error) {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return nil, fmt.Errorf("invalid order item ID format: %w", err)
	}

	var item models.OrderItem
	if err := s.DB.First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order item not found: %w", err)
		}
		return nil, fmt.Errorf("failed to get order item by ID: %w", err)
	}
	return &item, nil
}

func (s *Storage) UpdateOrderItem(itemID string, updates models.OrderItem) error {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return fmt.Errorf("invalid order item ID format: %w", err)
	}

	result := s.DB.Model(&models.OrderItem{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update order item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order item not found")
	}
	return nil
}

func (s *Storage) DeleteOrderItem(itemID string) error {
	id, err := uuid.Parse(itemID)
	if err != nil {
		return fmt.Errorf("invalid order item ID format: %w", err)
	}

	result := s.DB.Delete(&models.OrderItem{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete order item: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("order item not found")
	}
	return nil
}
