package service

import (
	"errors"
	"gorm.io/gorm"
	"task-manager/dto/request"
	"task-manager/dto/response"
	"task-manager/model"
	"task-manager/repository"
)

type TaskService interface {
	CreateTask(userID uint, req request.CreateTaskRequest) (*response.TaskResponse, error)
	GetTaskByID(id uint) (*response.TaskResponse, error)
	GetAllTasks(page, pageSize int) (*response.TaskListResponse, error)
	GetTasksByUserID(userID uint, page, pageSize int) (*response.TaskListResponse, error)
	UpdateTask(id uint, userID uint, req request.UpdateTaskRequest) (*response.TaskResponse, error)
	DeleteTask(id uint, userID uint) error
}

type taskService struct {
	taskRepo repository.TaskRepository
	userRepo repository.UserRepository
	db       *gorm.DB
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	userRepo repository.UserRepository,
	db *gorm.DB,
) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		userRepo: userRepo,
		db:       db,
	}
}

func (s *taskService) CreateTask(userID uint, req request.CreateTaskRequest) (*response.TaskResponse, error) {
	assignedUser, err := s.userRepo.FindByID(req.AssignedToID)
	if err != nil {
		return nil, err
	}
	if assignedUser == nil {
		return nil, errors.New("assigned user not found")
	}

	task := &model.Task{
		Title:        req.Title,
		Description:  req.Description,
		Status:       model.TaskStatusTodo,
		AssignedToID: req.AssignedToID,
		CreatedBy:    userID,
		DueDate:      req.DueDate,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		taskRepoTx := s.taskRepo.WithTx(tx)
		if err := taskRepoTx.Create(task); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	createdTask, err := s.taskRepo.FindByID(task.ID)
	if err != nil {
		return nil, err
	}

	taskResponse := response.MapTaskToResponse(*createdTask, true)
	return &taskResponse, nil
}

func (s *taskService) GetTaskByID(id uint) (*response.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	taskResponse := response.MapTaskToResponse(*task, true)
	return &taskResponse, nil
}

func (s *taskService) GetAllTasks(page, pageSize int) (*response.TaskListResponse, error) {
	offset := (page - 1) * pageSize
	tasks, total, err := s.taskRepo.FindAll(offset, pageSize)
	if err != nil {
		return nil, err
	}

	taskResponses := make([]response.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = response.MapTaskToResponse(task, true)
	}

	return &response.TaskListResponse{
		Tasks: taskResponses,
		Total: total,
	}, nil
}

func (s *taskService) GetTasksByUserID(userID uint, page, pageSize int) (*response.TaskListResponse, error) {
	offset := (page - 1) * pageSize
	tasks, total, err := s.taskRepo.FindByUserID(userID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	taskResponses := make([]response.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = response.MapTaskToResponse(task, true)
	}

	return &response.TaskListResponse{
		Tasks: taskResponses,
		Total: total,
	}, nil
}

func (s *taskService) UpdateTask(id uint, userID uint, req request.UpdateTaskRequest) (*response.TaskResponse, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = model.TaskStatus(*req.Status)
	}
	if req.AssignedToID != nil {
		assignedUser, err := s.userRepo.FindByID(*req.AssignedToID)
		if err != nil {
			return nil, err
		}
		if assignedUser == nil {
			return nil, errors.New("assigned user not found")
		}
		task.AssignedToID = *req.AssignedToID
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		taskRepoTx := s.taskRepo.WithTx(tx)
		if err := taskRepoTx.Update(task); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedTask, err := s.taskRepo.FindByID(task.ID)
	if err != nil {
		return nil, err
	}

	taskResponse := response.MapTaskToResponse(*updatedTask, true)
	return &taskResponse, nil
}

func (s *taskService) DeleteTask(id uint, userID uint) error {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		taskRepoTx := s.taskRepo.WithTx(tx)
		if err := taskRepoTx.Delete(id); err != nil {
			return err
		}
		return nil
	})
}