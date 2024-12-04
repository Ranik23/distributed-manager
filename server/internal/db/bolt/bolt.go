package bolt

import (
	"distributed-manager/server/internal/raft"
	"encoding/json"
	"fmt"
	"time"
	"github.com/boltdb/bolt"
)


type BoltDatabase struct {
	DB 		*bolt.DB
	LogPath string
}

// Конструктор для BoltDatabase
func NewBoltDatabase(logpath string) (*BoltDatabase, error) {

	db, err := bolt.Open(logpath, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}

	return &BoltDatabase{
		DB: db,
		LogPath: logpath,
	}, nil
}

// Метод для получения задачи по id
func (d *BoltDatabase) GetTask(id int) (raft.Task, error) {
	var task raft.Task

	// Открываем транзакцию для чтения
	err := d.DB.View(func(tx *bolt.Tx) error {
		// Ищем нужный бакет
		bucket := tx.Bucket([]byte("tasks"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		// Получаем значение по ключу
		taskData := bucket.Get(itob(id))
		if taskData == nil {
			return fmt.Errorf("task not found")
		}

		// Десериализуем данные
		if err := json.Unmarshal(taskData, &task); err != nil {
			return fmt.Errorf("error unmarshalling task data: %v", err)
		}

		return nil
	})

	if err != nil {
		return raft.Task{}, err
	}

	return task, nil
}

// Метод для отправки задачи в базу
func (d *BoltDatabase) SendTask(task raft.Task) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		if err != nil {
			return err
		}

		// Сериализуем задачу в JSON
		taskData, err := json.Marshal(task)
		if err != nil {
			return err
		}

		// Записываем данные в бакет
		return bucket.Put([]byte(task.ID), taskData)
	})
}

// Метод для удаления задачи по id
func (d *BoltDatabase) DeleteTask(id int) error {
	return d.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("tasks"))
		if bucket == nil {
			return fmt.Errorf("bucket not found")
		}

		// Удаляем задачу
		if err := bucket.Delete(itob(id)); err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		return nil
	})
}

// Функция для преобразования int в []byte
func itob(v int) []byte {
	return []byte(fmt.Sprintf("%d", v))
}