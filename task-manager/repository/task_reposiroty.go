package repository

import (
	"errors"
	"gorm.io/gorm"
	"task-manager/model"
)

type TaskRepository interface {
	Create(task *model.Task) error
	FindByID(id uint) (*model.Task, error)
	FindAll(offset, limit int) ([]model.Task, int64, error)
	FindByUserID(userID uint, offset, limit int) ([]model.Task, int64, error)
	Update(task *model.Task) error
	Delete(id uint) error
	WithTx(tx *gorm.DB) TaskRepository
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) WithTx(tx *gorm.DB) TaskRepository {
	return &taskRepository{db: tx}
}

func (r *taskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	if err := r.db.Preload("AssignedTo").First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) FindAll(offset, limit int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var count int64
	
	if err := r.db.Model(&model.Task{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	
	if err := r.db.Preload("AssignedTo").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	
	return tasks, count, nil
}

func (r *taskRepository) FindByUserID(userID uint, offset, limit int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var count int64
	
	if err := r.db.Model(&model.Task{}).Where("assigned_to_id = ?", userID).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	
	if err := r.db.Preload("AssignedTo").Where("assigned_to_id = ?", userID).
		Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}
	
	return tasks, count, nil
}

func (r *taskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&model.Task{}, id).Error
}