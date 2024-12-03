package storage

import (
	"distributed-manager/server/internal/entities/task"
	"log/slog"

	"gorm.io/gorm"
)

type Storage interface {
	GetTask(id int) (task.Task, error)
	SendTask(task task.Task) error
	DeleteTask(id int) error
	AutoMigrate(models interface{}) error	
}

type Database struct {
	db     *gorm.DB
	logger *slog.Logger
}

// NewDatabase создает новый экземпляр Database.
func NewDatabase(db *gorm.DB, logger *slog.Logger) *Database {
	return &Database{
		logger: logger,
		db:     db,
	}
}

// AutoMigrate выполняет автоматическую миграцию для указанных моделей.
func (d *Database) AutoMigrate(models interface{}) error {
	return d.db.AutoMigrate(models)
}

// GetTask получает задачу по ID.
func (d *Database) GetTask(id int) (task.Task, error) {
	var t task.Task
	result := d.db.First(&t, id) // Упрощенный запрос
	if result.Error != nil {
		d.logger.Error("failed to get task", slog.String("error", result.Error.Error()))
		return t, result.Error
	}
	d.logger.Info("got successfully", slog.Int("task_id", id))
	return t, nil
}

// SendTask отправляет задачу в базу данных.
func (d *Database) SendTask(t task.Task) error {
	result := d.db.Create(&t) // Упрощенный запрос
	if result.Error != nil {
		d.logger.Error("failed to send task", slog.String("error", result.Error.Error()))
		return result.Error
	}
	d.logger.Info("sent successfully", slog.String("task_id", string(t.ID))) // Предполагается, что у задачи есть поле ID
	return nil
}

// DeleteTask удаляет задачу по ID.
func (d *Database) DeleteTask(id int) error {
	result := d.db.Delete(&task.Task{}, id) // Исправленный запрос
	if result.Error != nil {
		d.logger.Error("failed to delete task", slog.String("error", result.Error.Error()))
		return result.Error
	}
	d.logger.Info("deleted successfully", slog.Int("task_id", id))
	return nil
}